package nebula

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/slackhq/nebula/cert"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
)

var curve = cert.Curve_CURVE25519

type CertInfo struct {
	Key  []byte
	Cert []byte
	Exp  time.Time
}

type SignReq struct {
	Name     string
	IP       string
	Groups   []string
	Subnets  []string
	Duration time.Duration
}

func CA(duration time.Duration) (*CertInfo, error) {
	pub, rawPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	exp := time.Now().Add(duration)
	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:      "mesh",
			NotBefore: time.Unix(time.Now().Unix(), 0),
			NotAfter:  time.Unix(exp.Unix(), 0),
			PublicKey: pub,
			IsCA:      true,
			Curve:     curve,
		},
	}

	if err := nc.Sign(curve, rawPriv); err != nil {
		return nil, err
	}

	pem, err := nc.MarshalToPEM()
	if err != nil {
		return nil, err
	}

	return &CertInfo{
		Key:  cert.MarshalSigningPrivateKey(curve, rawPriv),
		Cert: pem,
		Exp:  exp,
	}, nil
}

func Sign(rawCAKey, rawCACert []byte, req *SignReq) (*CertInfo, error) {
	caKey, _, _, err := cert.UnmarshalSigningPrivateKey(rawCAKey)
	if err != nil {
		return nil, err
	}

	caCert, _, err := cert.UnmarshalNebulaCertificateFromPEM(rawCACert)
	if err != nil {
		return nil, err
	}

	if err := caCert.VerifyPrivateKey(curve, caKey); err != nil {
		return nil, err
	}

	issuer, err := caCert.Sha256Sum()
	if err != nil {
		return nil, err
	}

	ip, ipNet, err := net.ParseCIDR(req.IP)
	if err != nil || ip.To4() == nil {
		return nil, fmt.Errorf("invalid ip definition: %s", req.IP)
	}
	ipNet.IP = ip

	if caCert.Expired(time.Now()) {
		return nil, fmt.Errorf("CA certificate has expired")
	}

	var subnets []*net.IPNet
	for _, subnet := range req.Subnets {
		_, s, err := net.ParseCIDR(subnet)
		if err != nil || s.IP.To4() == nil {
			return nil, fmt.Errorf("invalid subnet definition: %s", subnet)
		}
		subnets = append(subnets, s)
	}

	pub, rawPriv := x25519Keypair()
	exp := time.Now().Add(req.Duration)
	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:      req.Name,
			Ips:       []*net.IPNet{ipNet},
			Groups:    req.Groups,
			Subnets:   subnets,
			NotBefore: time.Unix(time.Now().Unix(), 0),
			NotAfter:  time.Unix(exp.Unix(), 0),
			PublicKey: pub,
			IsCA:      false,
			Issuer:    issuer,
			Curve:     caCert.Details.Curve,
		},
	}

	if err := nc.Sign(caCert.Details.Curve, caKey); err != nil {
		return nil, err
	}

	pem, err := nc.MarshalToPEM()
	if err != nil {
		return nil, err
	}

	return &CertInfo{
		Key:  cert.MarshalPrivateKey(curve, rawPriv),
		Cert: pem,
		Exp:  exp,
	}, nil
}

func x25519Keypair() ([]byte, []byte) {
	privkey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, privkey); err != nil {
		panic(err)
	}

	pubkey, err := curve25519.X25519(privkey, curve25519.Basepoint)
	if err != nil {
		panic(err)
	}

	return pubkey, privkey
}
