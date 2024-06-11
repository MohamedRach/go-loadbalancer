package main
/*
this is a simple project to refresh my golang knowledge
its a load blancer with round robin algorith

*/
import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type SimpleServer struct {
  addr string
  proxy *httputil.ReverseProxy
}

type LoadBalancer struct {
  port string
  roundRobinCount int
  servers []Server

}
func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
  return &LoadBalancer{
    port: port,
    roundRobinCount: 0,
    servers: servers,
  }

}

func (lb *LoadBalancer) getNextAvailableServer()  Server{
  server := lb.servers[lb.roundRobinCount%len(lb.servers)]
  if !server.IsAlive() {
    lb.roundRobinCount++
    server = lb.servers[lb.roundRobinCount%len(lb.servers)]
  }

  lb.roundRobinCount++
  return server
}

func (lb *LoadBalancer) serverProxy(rw http.ResponseWriter, r *http.Request)  {
  target_server := lb.getNextAvailableServer()
  fmt.Printf("forwarding requests to address %q\n", target_server.Address())
  target_server.Serve(rw, r)
}


type Server interface {
  Address() string
  IsAlive() bool
  Serve(rw http.ResponseWriter, r *http.Request)

}

func (s *SimpleServer) Address() string {return s.addr}

func (s *SimpleServer) IsAlive() bool { return true}

func (s *SimpleServer) Serve(rw http.ResponseWriter, r *http.Request){
  s.proxy.ServeHTTP(rw, r)
}




func NewSimpleServer(addr string) *SimpleServer  {
  serverUrl, err := url.Parse(addr)  
  handleErr(err)
  
  return &SimpleServer{
    addr: addr,
    proxy: httputil.NewSingleHostReverseProxy(serverUrl),
  }
}






func handleErr(err error) {
  if err != nil {
    fmt.Println("error: %v\n", err)

  }
  
}


func main() {
  servers := []Server {
    NewSimpleServer("https://www.facebook.com/"),
    NewSimpleServer("https://www.google.com/"),
    NewSimpleServer("https://www.youtube.com/"),
  }
  lb := NewLoadBalancer("8000", servers)
  handleRedirect := func(rw http.ResponseWriter, r *http.Request)  {
    lb.serverProxy(rw, r)
  }

  http.HandleFunc("/", handleRedirect)
  fmt.Printf("Serving requests at 'localhost: %s'\n", lb.port)
  http.ListenAndServe(":"+lb.port, nil)

}

  








