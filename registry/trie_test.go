package registry

import (
	"fmt"
	"testing"
)

func TestWordTrie1(t *testing.T) {
	trie := newWordTrie()
	trie.add("/")
	trie.add("/a/b")
	trie.add("/a/b/c")
	trie.delete("/a/b")
	trie.delete("/a/b/c")
	trie.debugTrace()
}

func TestWordTrie2(t *testing.T) {
	trie := newWordTrie()
	trie.add("/")
	trie.add("/a/b")
	trie.add("/a/c")
	trie.add("/a/c/d1")
	trie.add("/a/c/d2")
	trie.delete("/a/c/d2")
	trie.delete("/a")
	trie.delete("/a/c")
	trie.delete("/a/c/d1")
	trie.delete("/a/b")
	trie.delete("/")

	trie.debugTrace()
}

func TestWordTrie3(t *testing.T) {
	trie := newWordTrie()
	trie.add("/")
	trie.add("/a/b")
	trie.add("/a/c")
	trie.add("/a/c/d1")
	trie.add("/a/c/d2")
	fmt.Println("search /", trie.search("/a/c/"))

	trie.debugTrace()
}

func TestWordTrie4(t *testing.T) {
	trie := newWordTrie()
	trie.add("/a")
	trie.add("/a/c")
	trie.add("/b")

	fmt.Println("search /", trie.search("/"))

	trie.debugTrace()
}

func TestWordTrie5(t *testing.T) {
	trie := newWordTrie()
	trie.add("/a")
	trie.add("/a/c")
	trie.add("/a/b")

	fmt.Println("search /", trie.search("/a/"))

	trie.debugTrace()
}

func TestWordTrie6(t *testing.T) {
	trie := newWordTrie()
	trie.add("admin_grpc.Admin/admin_grpc.Admin_4001")
	trie.add("peer_grpc.Peer/peer_grpc.Peer_4001")
	trie.add("admin_grpc.Admin/admin_grpc.Admin_1001")
	trie.delete("peer_grpc.Peer/peer_grpc.Peer_4001")

	fmt.Println("search /", trie.search("admin_grpc.Admin/"))

	trie.debugTrace()
}
