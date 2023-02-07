package libWorker

import (
	"context"
	"errors"
	"math"
	"net/url"
	"sync"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/dop251/goja"
	"golang.org/x/time/rate"
)

const unverifiedDomainLimit = 10

type DomainLimiter struct {
	skipVerifiedDomainCheck bool
	standalone              bool
	limiters                map[string]*rate.Limiter
	limitersMutex           sync.Mutex
	verifiedDomains         []string
	workerInfo              *WorkerInfo
}

func CreateDomainLimiter(standalone bool, verifiedDomains []string, workerInfo *WorkerInfo) *DomainLimiter {
	return &DomainLimiter{
		skipVerifiedDomainCheck: lib.GetEnvVariableBool("SKIP_VERIFIED_DOMAIN_CHECK", false),
		standalone:              standalone,
		limiters:                make(map[string]*rate.Limiter),
		verifiedDomains:         verifiedDomains,
		workerInfo:              workerInfo,
	}
}

func getDomainFromURL(url interface{}) (string, error) {
	if urlJSValue, ok := url.(goja.Value); ok {
		url = urlJSValue.Export()
	}
	u, err := toURL(url)
	if err != nil {
		return "", err
	}

	domain := domainutil.Domain(u.getURL().Hostname())

	if domain == "" {
		// Try and get the ip instead
		ip := u.getURL().Hostname()

		if ip == "" {
			return "", errors.New("could not extract domain from url")
		}

		return ip, nil
	}

	return domain, nil
}

func (dl *DomainLimiter) WaitForLimiterGojaURL(url interface{}) error {
	domain, err := getDomainFromURL(url)
	if err != nil {
		return err
	}

	return dl.WaitForLimiter(domain)
}

func (dl *DomainLimiter) WaitForLimiterURL(url url.URL) error {
	domain := domainutil.Domain(url.Hostname())

	if domain == "" {
		// Try and get the ip instead
		ip := url.Hostname()

		if ip == "" {
			return errors.New("could not extract domain from url")
		}

		return dl.WaitForLimiter(ip)
	}

	return dl.WaitForLimiter(domain)
}

func (dl *DomainLimiter) WaitForLimiter(domain string) error {
	if !dl.standalone || dl.skipVerifiedDomainCheck {
		return nil
	}

	dl.limitersMutex.Lock()

	if _, ok := dl.limiters[domain]; !ok {
		dl.createDomainLimiter(domain)
	}

	limiter := dl.limiters[domain]

	dl.limitersMutex.Unlock()

	// Wait for the limiter to allow the request
	if err := limiter.Wait(context.Background()); err != nil {
		return err
	}

	return nil
}

func (dl *DomainLimiter) createDomainLimiter(domain string) {
	// Check if domain is in verified domains
	verified := false

	for _, verifiedDomain := range dl.verifiedDomains {
		if verifiedDomain == domain {
			verified = true
			break
		}
	}

	limit := math.MaxFloat64

	if !verified {
		DispatchMessage(*dl.workerInfo.Gs, "UNVERIFIED_DOMAIN_THROTTLED", "MESSAGE")
		limit = unverifiedDomainLimit * dl.workerInfo.SubFraction
	}

	dl.limiters[domain] = rate.NewLimiter(rate.Limit(limit), 1)
}

type urlInstance struct {
	u           *url.URL
	Name        string // http://example.com/thing/${}/
	urlInstance string // http://example.com/thing/1234/
	CleanURL    string // urlInstance with masked user credentials, used for output
}

func newURL(urlString, name string) (urlInstance, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return urlInstance{}, err
	}
	newURL := urlInstance{u: u, Name: name, urlInstance: urlString}
	newURL.CleanURL = newURL.clean()
	if urlString == name {
		newURL.Name = newURL.CleanURL
	}
	return newURL, nil
}

func (u urlInstance) clean() string {
	if u.CleanURL != "" {
		return u.CleanURL
	}

	if u.u == nil || u.u.User == nil {
		return u.urlInstance
	}

	if password, passwordOk := u.u.User.Password(); passwordOk {
		// here 3 is for the '://' and 4 is because of '://' and ':' between the credentials
		return u.urlInstance[:len(u.u.Scheme)+3] + "****:****" + u.urlInstance[len(u.u.Scheme)+4+len(u.u.User.Username())+len(password):]
	}

	// here 3 in both places is for the '://'
	return u.urlInstance[:len(u.u.Scheme)+3] + "****" + u.urlInstance[len(u.u.Scheme)+3+len(u.u.User.Username()):]
}

func (u urlInstance) getURL() *url.URL {
	return u.u
}

func toURL(u interface{}) (urlInstance, error) {
	switch tu := u.(type) {
	case urlInstance:
		// Handling of http.url`http://example.com/{$id}`
		return tu, nil
	case string:
		// Handling of "http://example.com/"
		return newURL(tu, tu)
	default:
		return urlInstance{}, errors.New("invalid url")
	}
}
