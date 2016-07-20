package user_ip

import (
  "net/http"
  "net"
  "strings"
  "unicode"
  "bytes"
  "github.com/kataras/iris"
)

type PrivateAddressesRange struct {
  start net.IP
  end net.IP
}

var privateIpAddresses = []PrivateAddressesRange{
  PrivateAddressesRange{
    start: net.IPv4(10, 0, 0, 0),
    end: net.IPv4(10, 255, 255, 255),
  },
  PrivateAddressesRange{
    start: net.IPv4(100, 64, 0, 0),
    end: net.IPv4(100, 127, 255, 255),
  },
  PrivateAddressesRange{
    start: net.IPv4(172, 16, 0, 0),
    end: net.IPv4(172, 31, 255, 255),
  },
  PrivateAddressesRange{
    start: net.IPv4(192, 0, 0, 0),
    end: net.IPv4(192, 0, 0, 255),
  },
  PrivateAddressesRange{
    start: net.IPv4(192, 168, 0, 0),
    end: net.IPv4(192, 168, 255, 255),
  },
  PrivateAddressesRange{
    start: net.IPv4(198, 18, 0, 0),
    end: net.IPv4(198, 19, 255, 255),
  },
  PrivateAddressesRange {
    start: net.ParseIP("fc00::"),
    end: net.ParseIP("fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
  },
}

var HEADER_IP_KEYS = []string{"X-Forwarded-For", "X-Real-Ip"}

const MAX_IP_SEARCH int = 4

// =====================================================================================

func removeWhiteSpaces(str string) string {
  return strings.Map(func(char rune) rune {
    if unicode.IsSpace(char) {
      return -1
    }
    return char
  }, str)
}

func isPrivateAddress(ip net.IP) bool {
  for _, privateIpRange := range privateIpAddresses {
    if bytes.Compare(ip, privateIpRange.start) >= 0 && bytes.Compare(ip, privateIpRange.end) <= 0 {
      return true
    }
  }

  return false
}

func ignoreIp(ip net.IP) bool {
  return !ip.IsGlobalUnicast() || isPrivateAddress(ip)
}

func findRealIP(ipAddress *string, availableIps []string) {
  if len(availableIps) >= MAX_IP_SEARCH {
    availableIps = availableIps[:MAX_IP_SEARCH]
  }

  for _, ip := range availableIps {
    ip = removeWhiteSpaces(ip)
    userIp := net.ParseIP(ip)

    if ignoreIp(userIp) {
      continue
    } else {
      *ipAddress = ip
    }
  }
}

// =====================================================================================
func GetFromContext(ctx *iris.Context) string {
  var ipAddress string

  for _, headerProperty := range HEADER_IP_KEYS {
    availableIps := strings.Split(string(ctx.RequestCtx.Request.Header.Peek(headerProperty)), ",")
    findRealIP(&ipAddress, availableIps)
  }

  return ipAddress
}

func Get(req *http.Request) string {
  var ipAddress string

  for _, headerProperty := range HEADER_IP_KEYS {

    // We will only search through up to MAX_IP_SEARCH, we should not be looking
    // for all of them, since client may be playing with X-Forwarded-For
    availableIps := strings.Split(req.Header.Get(headerProperty), ",")
    findRealIP(&ipAddress, availableIps)
  }

  return ipAddress
}
