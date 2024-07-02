package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

var rootServers = []string{
	"198.41.0.4",     // a.root-servers.net
	"199.9.14.201",   // b.root-servers.net
	"192.33.4.12",    // c.root-servers.net
	"199.7.91.13",    // d.root-servers.net
	"192.203.230.10", // e.root-servers.net
	"192.5.5.241",    // f.root-servers.net
	"192.112.36.4",   // g.root-servers.net
	"198.97.190.53",  // h.root-servers.net
	"192.36.148.17",  // i.root-servers.net
	"192.58.128.30",  // j.root-servers.net
	"193.0.14.129",   // k.root-servers.net
	"199.7.83.42",    // l.root-servers.net
	"202.12.27.33",   // m.root-servers.net
}

func main() {
	server := &dns.Server{Addr: ":8153", Net: "udp"}
	dns.HandleFunc(".", handleDNSRequest)
	log.Printf("Starting DNS server on port 8153")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	for _, question := range r.Question {
		log.Printf("Received query for %s\n", question.Name)
		response, err := recursiveResolve(rootServers, r.Question[0].Name, dns.TypeA)

		if err != nil {
			log.Printf("Failed to resolve %s: %s\n", question.Name, err)
			dns.HandleFailed(w, r)
			return
		}
		m.Answer = append(m.Answer, response.Answer...)
	}
	w.WriteMsg(m)
}

func recursiveResolve(servers []string, name string, qtype uint16) (*dns.Msg, error) {
	r := new(dns.Msg)
	r.SetQuestion(name, qtype)

	fmt.Println("SERVERS:", servers)

	new_servers := make([]string, 0)

	servers_map := make(map[string]int, 0)

	for _, rootServer := range servers {
		fmt.Println(rootServer)
		response, err := queryServer(rootServer, r)
		if err != nil {
			continue
		}

		if len(response.Answer) > 0 {
			for _, rr := range response.Answer {
				if _, ok := rr.(*dns.A); ok {
					return response, nil
				}
			}
		}

		namespaces_map := make(map[string]int, 0)
		for _, rr := range response.Ns {
			if nsRR, ok := rr.(*dns.NS); ok {
				if nsRR.Hdr.Rrtype == dns.TypeNS {
					fmt.Println("Namespace: ", nsRR)
					namespaces_map[nsRR.Ns] = 1
				}
			}

		}

		for _, rr := range response.Extra {
			if aRR, ok := rr.(*dns.A); ok {
				namespace := aRR.Hdr.Name
				_, ns_exists := namespaces_map[namespace]
				_, server_exists := servers_map[aRR.A.String()]
				if ns_exists && !server_exists {
					new_servers = append(new_servers, aRR.A.String())
					servers_map[aRR.A.String()] = 1
				}
			}
		}
	}

	if len(new_servers) > 0 {
		return recursiveResolve(new_servers, name, qtype)
	}
	return nil, log.Output(2, "All root server queries failed")
}

func queryServer(server string, r *dns.Msg) (*dns.Msg, error) {
	c := new(dns.Client)
	response, _, err := c.Exchange(r, net.JoinHostPort(server, "53"))
	return response, err
}
