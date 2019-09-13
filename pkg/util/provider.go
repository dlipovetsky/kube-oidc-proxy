package util

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"time"

	oidc "github.com/coreos/go-oidc"
	netutil "k8s.io/apimachinery/pkg/util/net"
	certutil "k8s.io/client-go/util/cert"

	"github.com/jetstack/kube-oidc-proxy/cmd/options"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
)

func InitProvider(ctx context.Context, opts *options.OIDCAuthenticationOptions, stopCh <-chan struct{}) error {
	url, err := url.Parse(opts.IssuerURL)
	if err != nil {
		return err
	}

	if url.Scheme != "https" {
		return fmt.Errorf("'oidc-issuer-url' (%q) has invalid scheme (%q), require 'https'", opts.IssuerURL, url.Scheme)
	}

	var roots *x509.CertPool
	if opts.CAFile != "" {
		roots, err = certutil.NewPool(opts.CAFile)
		if err != nil {
			return fmt.Errorf("Failed to read the CA file: %v", err)
		}
	} else {
		klog.Info("OIDC: No x509 certificates provided, will use host's root CA set")
	}

	// Copied from http.DefaultTransport.
	tr := netutil.SetTransportDefaults(&http.Transport{
		// According to golang's doc, if RootCAs is nil,
		// TLS uses the host's root CA set.
		TLSClientConfig: &tls.Config{RootCAs: roots},
	})

	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	ctx = oidc.ClientContext(ctx, client)
	return wait.PollUntil(time.Second*10, func() (bool, error) {
		_, err := oidc.NewProvider(ctx, opts.IssuerURL)
		if err != nil {
			klog.Errorf("failed to initialize oidc provider: %v", err)
			return false, nil
		}
		return true, nil
	}, stopCh)
}
