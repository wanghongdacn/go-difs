/*
Package core implements the IpfsNode object and related methods.

Packages underneath core/ provide a (relatively) stable, low-level API
to carry out most IPFS-related tasks.  For more details on the other
interfaces and how core/... fits into the bigger IPFS picture, see:

  $ godoc github.com/daccproject/go-difs
*/
package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	version "github.com/daccproject/go-difs"
	rp "github.com/daccproject/go-difs/exchange/reprovide"
	filestore "github.com/daccproject/go-difs/filestore"
	mount "github.com/daccproject/go-difs/fuse/mount"
	iam "github.com/daccproject/go-difs/iam"
	namesys "github.com/daccproject/go-difs/namesys"
	ipnsrp "github.com/daccproject/go-difs/namesys/republisher"
	p2p "github.com/daccproject/go-difs/p2p"
	pin "github.com/daccproject/go-difs/pin"
	repo "github.com/daccproject/go-difs/repo"

	u "gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	psrouter "gx/ipfs/QmPksYjuFjt2M63ufuDYEFdwVAmFcDGJgoED8T8dwfcwED/go-libp2p-pubsub-router"
	ic "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	libp2p "gx/ipfs/QmQiaskfWpdRJ4x2spEQjPFTUkEB87KDYu91qnNYBqvvcX/go-libp2p"
	discovery "gx/ipfs/QmQiaskfWpdRJ4x2spEQjPFTUkEB87KDYu91qnNYBqvvcX/go-libp2p/p2p/discovery"
	p2pbhost "gx/ipfs/QmQiaskfWpdRJ4x2spEQjPFTUkEB87KDYu91qnNYBqvvcX/go-libp2p/p2p/host/basic"
	rhost "gx/ipfs/QmQiaskfWpdRJ4x2spEQjPFTUkEB87KDYu91qnNYBqvvcX/go-libp2p/p2p/host/routed"
	identify "gx/ipfs/QmQiaskfWpdRJ4x2spEQjPFTUkEB87KDYu91qnNYBqvvcX/go-libp2p/p2p/protocol/identify"
	ping "gx/ipfs/QmQiaskfWpdRJ4x2spEQjPFTUkEB87KDYu91qnNYBqvvcX/go-libp2p/p2p/protocol/ping"
	ft "gx/ipfs/QmQjEpRiwVvtowhq69dAtB4jhioPVFXiCcWZm9Sfgn7eqc/go-unixfs"
	peer "gx/ipfs/QmQsErDt8Qgw1XrsXf2BpEzDgGWtB1YLsTAARBup5b6B9W/go-libp2p-peer"
	connmgr "gx/ipfs/QmRAvQVfAFfzqwk1VBYJMci9SrqJcMULNJYeZzzprDRn5V/go-libp2p-connmgr"
	p2phost "gx/ipfs/QmRRCrNRs4qxotXx7WJT6SpCvSNEhXvyBcVjXY2K71pcjE/go-libp2p-host"
	logging "gx/ipfs/QmRREK2CAZ5Re2Bd9zZFG6FeYDppUWt5cMgsoUEp3ktgSr/go-log"
	merkledag "gx/ipfs/QmRiQCJZ91B7VNmLvA6sxzDuBJGSojS3uXHHVuNr3iueNZ/go-merkledag"
	routing "gx/ipfs/QmS4niovD1U6pRjUBXivr1zvvLBqiTKbERjFo994JU7oQS/go-libp2p-routing"
	goprocess "gx/ipfs/QmSF8fPo3jgVBAy8fpdjjYqgG87dkJgUprRBHRd2tmfgpP/goprocess"
	mamask "gx/ipfs/QmSMZwvs3n4GBikZ7hKzT17c3bk65FmyZo2JqtJ16swqCv/multiaddr-filter"
	mafilter "gx/ipfs/QmSW4uNHbvQia8iZDXzbwjiyHQtnyo9aFqfQAMasj3TJ6Y/go-maddr-filter"
	dht "gx/ipfs/QmTRj8mj6X5LtjVochPPSNX6MTbJ6iVojcfakWJKG13re7/go-libp2p-kad-dht"
	dhtopts "gx/ipfs/QmTRj8mj6X5LtjVochPPSNX6MTbJ6iVojcfakWJKG13re7/go-libp2p-kad-dht/opts"
	bitswap "gx/ipfs/QmTtmrK4iiM3MxWNA3pvbM9ekQiGZAiFyo57GP8B9FFgtz/go-bitswap"
	bsnet "gx/ipfs/QmTtmrK4iiM3MxWNA3pvbM9ekQiGZAiFyo57GP8B9FFgtz/go-bitswap/network"
	config "gx/ipfs/QmTyiSs9VgdVb4pnzdjtKhcfdTkHFEaNn6xnCbZq4DTFRt/go-ipfs-config"
	ifconnmgr "gx/ipfs/QmUPz6FCzCCU7sTY9Sore5NGSUA8YSF2yMkLPjDFq7wGqD/go-libp2p-interface-connmgr"
	floodsub "gx/ipfs/QmVFB6rGJEZnzJrQwoEhbyDs1tA8RVsQvCS6JXpuw9Xtta/go-libp2p-floodsub"
	ds "gx/ipfs/QmVG5gxteQNEMhrS8prJSmU2C9rebtFuTd3SYZ5kE3YZ5k/go-datastore"
	exchange "gx/ipfs/QmWw71Mz9PXKgYG8ZfTYN7Ax2Zm48Eurbne3wC2y7CKmLz/go-ipfs-exchange-interface"
	ipld "gx/ipfs/QmX5CsuHyVZeTLxgRSYkgLSDQKb9UjE8xnhQzCEJWWWFsC/go-ipld-format"
	smux "gx/ipfs/QmY9JXR3FupnYAYJWK9aMr9bCpqWKcToQ1tz8DVGTrHpHw/go-stream-muxer"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	cid "gx/ipfs/QmZFbDTY9jfSBms2MchvYM9oYRbAF19K7Pby47yDBfpPrb/go-cid"
	pnet "gx/ipfs/QmZaQ3K9PRd5sYYoG1xbTGPtd3N7TYiKBRmcBUTsx8HVET/go-libp2p-pnet"
	bserv "gx/ipfs/QmbSB9Uh3wVgmiCb1fAb8zuC3qAE6un4kd1jvatUurfAmB/go-blockservice"
	circuit "gx/ipfs/Qmbc6WjgbkaYhPw5dd6X2RLGiJz854dPPKkyZf52vE2PTS/go-libp2p-circuit"
	bstore "gx/ipfs/QmcmpX42gtDv1fz24kau4wjS9hfwWj5VexWBKgGnWzsyag/go-ipfs-blockstore"
	yamux "gx/ipfs/QmcsgrV3nCAKjiHKZhKVXWc4oY3WBECJCqahXEMpHeMrev/go-smux-yamux"
	rhelpers "gx/ipfs/Qmd22J9AnyR3QUH56WPXkrTbCNkQ4x7TWWinHcZBhQkgVw/go-libp2p-routing-helpers"
	nilrouting "gx/ipfs/Qmd45r5jHr1PKMNQqifnbZy1ZQwHdtXUDJFamUEvUJE544/go-ipfs-routing/none"
	offroute "gx/ipfs/Qmd45r5jHr1PKMNQqifnbZy1ZQwHdtXUDJFamUEvUJE544/go-ipfs-routing/offline"
	record "gx/ipfs/QmdHb9aBELnQKTVhvvA3hsQbRgUAwsWUzBP2vZ6Y5FBYvE/go-libp2p-record"
	"gx/ipfs/QmdMPBephdLYNESkruDX2hcDTgFYhoCt4LimWhgnomSdV2/go-path/resolver"
	mfs "gx/ipfs/QmdghKsSDa2AD1kC4qYRnVYWqZecdSBRZjeXRdhMYYhafj/go-mfs"
	metrics "gx/ipfs/QmdhwKw53CTV8EJSAsR1bpmMT5kXiWBgeAyv1EXeeDiXqR/go-libp2p-metrics"
	mplex "gx/ipfs/QmdiBZzwGtN2yHJrWD9ojQ7ASS48nv7BcojWLkYd1ZtrV2/go-smux-multiplex"
	pstore "gx/ipfs/QmeKD8YT7887Xu6Z86iZmpYNxrLogJexqxEugSmaf14k64/go-libp2p-peerstore"
)

const IpnsValidatorTag = "ipns"

const kReprovideFrequency = time.Hour * 12
const discoveryConnTimeout = time.Second * 30

var log = logging.Logger("core")

type mode int

const (
	// zero value is not a valid mode, must be explicitly set
	localMode mode = iota
	offlineMode
	onlineMode
)

func init() {
	identify.ClientVersion = "go-ipfs/" + version.CurrentVersionNumber + "/" + version.CurrentCommit
}

// IpfsNode is IPFS Core module. It represents an IPFS instance.
type IpfsNode struct {

	// Self
	Identity peer.ID // the local node's identity

	Repo repo.Repo

	// Local node
	Pinning         pin.Pinner // the pinning manager
	Mounts          Mounts     // current mount state, if any.
	PrivateKey      ic.PrivKey // the local node's private Key
	PNetFingerprint []byte     // fingerprint of private network

	// Services
	Peerstore       pstore.Peerstore     // storage for other Peer instances
	Blockstore      bstore.GCBlockstore  // the block store (lower level)
	Filestore       *filestore.Filestore // the filestore blockstore
	BaseBlocks      bstore.Blockstore    // the raw blockstore, no filestore wrapping
	GCLocker        bstore.GCLocker      // the locker used to protect the blockstore during gc
	Blocks          bserv.BlockService   // the block service, get/add blocks.
	DAG             ipld.DAGService      // the merkle dag service, get/add objects.
	Resolver        *resolver.Resolver   // the path resolution system
	Reporter        metrics.Reporter
	Discovery       discovery.Service
	FilesRoot       *mfs.Root
	RecordValidator record.Validator

	// Online
	PeerHost     p2phost.Host        // the network host (server+client)
	Bootstrapper io.Closer           // the periodic bootstrapper
	Routing      routing.IpfsRouting // the routing system. recommend ipfs-dht
	Exchange     exchange.Interface  // the block exchange + strategy (bitswap)
	Namesys      namesys.NameSystem  // the name system, resolves paths to hashes
	Ping         *ping.PingService
	Reprovider   *rp.Reprovider // the value reprovider system
	IpnsRepub    *ipnsrp.Republisher

	Floodsub *floodsub.PubSub
	PSRouter *psrouter.PubsubValueStore
	DHT      *dht.IpfsDHT
	IAM      *iam.IAM
	P2P      *p2p.P2P

	proc goprocess.Process
	ctx  context.Context

	mode         mode
	localModeSet bool
}

// Mounts defines what the node's mount state is. This should
// perhaps be moved to the daemon or mount. It's here because
// it needs to be accessible across daemon requests.
type Mounts struct {
	Ipfs mount.Mount
	Ipns mount.Mount
}

func (n *IpfsNode) startOnlineServices(ctx context.Context, routingOption RoutingOption, hostOption HostOption, do DiscoveryOption, pubsub, ipnsps, mplex bool, iamOption IAMOption) error {
	if n.PeerHost != nil { // already online.
		return errors.New("node already online")
	}

	// load private key
	if err := n.LoadPrivateKey(); err != nil {
		return err
	}

	// get undialable addrs from config
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	var libp2pOpts []libp2p.Option
	for _, s := range cfg.Swarm.AddrFilters {
		f, err := mamask.NewMask(s)
		if err != nil {
			return fmt.Errorf("incorrectly formatted address filter in config: %s", s)
		}
		libp2pOpts = append(libp2pOpts, libp2p.FilterAddresses(f))
	}

	if !cfg.Swarm.DisableBandwidthMetrics {
		// Set reporter
		n.Reporter = metrics.NewBandwidthCounter()
		libp2pOpts = append(libp2pOpts, libp2p.BandwidthReporter(n.Reporter))
	}

	swarmkey, err := n.Repo.SwarmKey()
	if err != nil {
		return err
	}

	if swarmkey != nil {
		protec, err := pnet.NewProtector(bytes.NewReader(swarmkey))
		if err != nil {
			return fmt.Errorf("failed to configure private network: %s", err)
		}
		n.PNetFingerprint = protec.Fingerprint()
		go func() {
			t := time.NewTicker(30 * time.Second)
			<-t.C // swallow one tick
			for {
				select {
				case <-t.C:
					if ph := n.PeerHost; ph != nil {
						if len(ph.Network().Peers()) == 0 {
							log.Warning("We are in private network and have no peers.")
							log.Warning("This might be configuration mistake.")
						}
					}
				case <-n.Process().Closing():
					t.Stop()
					return
				}
			}
		}()

		libp2pOpts = append(libp2pOpts, libp2p.PrivateNetwork(protec))
	}

	addrsFactory, err := makeAddrsFactory(cfg.Addresses)
	if err != nil {
		return err
	}
	if !cfg.Swarm.DisableRelay {
		addrsFactory = composeAddrsFactory(addrsFactory, filterRelayAddrs)
	}
	libp2pOpts = append(libp2pOpts, libp2p.AddrsFactory(addrsFactory))

	connm, err := constructConnMgr(cfg.Swarm.ConnMgr)
	if err != nil {
		return err
	}
	libp2pOpts = append(libp2pOpts, libp2p.ConnectionManager(connm))

	libp2pOpts = append(libp2pOpts, makeSmuxTransportOption(mplex))

	if !cfg.Swarm.DisableNatPortMap {
		libp2pOpts = append(libp2pOpts, libp2p.NATPortMap())
	}
	if !cfg.Swarm.DisableRelay {
		var opts []circuit.RelayOpt
		if cfg.Swarm.EnableRelayHop {
			opts = append(opts, circuit.OptHop)
		}
		libp2pOpts = append(libp2pOpts, libp2p.EnableRelay(opts...))
	}

	peerhost, err := hostOption(ctx, n.Identity, n.Peerstore, libp2pOpts...)

	if err != nil {
		return err
	}

	if err := n.startOnlineServicesWithHost(ctx, peerhost, routingOption, pubsub, ipnsps, iamOption); err != nil {
		return err
	}

	// Ok, now we're ready to listen.
	if err := startListening(n.PeerHost, cfg); err != nil {
		return err
	}

	n.P2P = p2p.NewP2P(n.Identity, n.PeerHost, n.Peerstore)

	// setup local discovery
	if do != nil {
		service, err := do(ctx, n.PeerHost)
		if err != nil {
			log.Error("mdns error: ", err)
		} else {
			service.RegisterNotifee(n)
			n.Discovery = service
		}
	}

	return n.Bootstrap(DefaultBootstrapConfig)
}

func constructConnMgr(cfg config.ConnMgr) (ifconnmgr.ConnManager, error) {
	switch cfg.Type {
	case "":
		// 'default' value is the basic connection manager
		return connmgr.NewConnManager(config.DefaultConnMgrLowWater, config.DefaultConnMgrHighWater, config.DefaultConnMgrGracePeriod), nil
	case "none":
		return nil, nil
	case "basic":
		grace, err := time.ParseDuration(cfg.GracePeriod)
		if err != nil {
			return nil, fmt.Errorf("parsing Swarm.ConnMgr.GracePeriod: %s", err)
		}

		return connmgr.NewConnManager(cfg.LowWater, cfg.HighWater, grace), nil
	default:
		return nil, fmt.Errorf("unrecognized ConnMgr.Type: %q", cfg.Type)
	}
}

func (n *IpfsNode) startLateOnlineServices(ctx context.Context) error {
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	var keyProvider rp.KeyChanFunc

	switch cfg.Reprovider.Strategy {
	case "all":
		fallthrough
	case "":
		keyProvider = rp.NewBlockstoreProvider(n.Blockstore)
	case "roots":
		keyProvider = rp.NewPinnedProvider(n.Pinning, n.DAG, true)
	case "pinned":
		keyProvider = rp.NewPinnedProvider(n.Pinning, n.DAG, false)
	default:
		return fmt.Errorf("unknown reprovider strategy '%s'", cfg.Reprovider.Strategy)
	}
	n.Reprovider = rp.NewReprovider(ctx, n.Routing, keyProvider)

	reproviderInterval := kReprovideFrequency
	if cfg.Reprovider.Interval != "" {
		dur, err := time.ParseDuration(cfg.Reprovider.Interval)
		if err != nil {
			return err
		}

		reproviderInterval = dur
	}

	go n.Reprovider.Run(reproviderInterval)

	return nil
}

func makeAddrsFactory(cfg config.Addresses) (p2pbhost.AddrsFactory, error) {
	var annAddrs []ma.Multiaddr
	for _, addr := range cfg.Announce {
		maddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, err
		}
		annAddrs = append(annAddrs, maddr)
	}

	filters := mafilter.NewFilters()
	noAnnAddrs := map[string]bool{}
	for _, addr := range cfg.NoAnnounce {
		f, err := mamask.NewMask(addr)
		if err == nil {
			filters.AddDialFilter(f)
			continue
		}
		maddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, err
		}
		noAnnAddrs[maddr.String()] = true
	}

	return func(allAddrs []ma.Multiaddr) []ma.Multiaddr {
		var addrs []ma.Multiaddr
		if len(annAddrs) > 0 {
			addrs = annAddrs
		} else {
			addrs = allAddrs
		}

		var out []ma.Multiaddr
		for _, maddr := range addrs {
			// check for exact matches
			ok, _ := noAnnAddrs[maddr.String()]
			// check for /ipcidr matches
			if !ok && !filters.AddrBlocked(maddr) {
				out = append(out, maddr)
			}
		}
		return out
	}, nil
}

func makeSmuxTransportOption(mplexExp bool) libp2p.Option {
	const yamuxID = "/yamux/1.0.0"
	const mplexID = "/mplex/6.7.0"

	ymxtpt := &yamux.Transport{
		AcceptBacklog:          512,
		ConnectionWriteTimeout: time.Second * 10,
		KeepAliveInterval:      time.Second * 30,
		EnableKeepAlive:        true,
		MaxStreamWindowSize:    uint32(1024 * 512),
		LogOutput:              ioutil.Discard,
	}

	if os.Getenv("YAMUX_DEBUG") != "" {
		ymxtpt.LogOutput = os.Stderr
	}

	muxers := map[string]smux.Transport{yamuxID: ymxtpt}
	if mplexExp {
		muxers[mplexID] = mplex.DefaultTransport
	}

	// Allow muxer preference order overriding
	order := []string{yamuxID, mplexID}
	if prefs := os.Getenv("LIBP2P_MUX_PREFS"); prefs != "" {
		order = strings.Fields(prefs)
	}

	opts := make([]libp2p.Option, 0, len(order))
	for _, id := range order {
		tpt, ok := muxers[id]
		if !ok {
			log.Warning("unknown or duplicate muxer in LIBP2P_MUX_PREFS: %s", id)
			continue
		}
		delete(muxers, id)
		opts = append(opts, libp2p.Muxer(id, tpt))
	}

	return libp2p.ChainOptions(opts...)
}

func setupDiscoveryOption(d config.Discovery) DiscoveryOption {
	if d.MDNS.Enabled {
		return func(ctx context.Context, h p2phost.Host) (discovery.Service, error) {
			if d.MDNS.Interval == 0 {
				d.MDNS.Interval = 5
			}
			return discovery.NewMdnsService(ctx, h, time.Duration(d.MDNS.Interval)*time.Second, discovery.ServiceTag)
		}
	}
	return nil
}

// HandlePeerFound attempts to connect to peer from `PeerInfo`, if it fails
// logs a warning log.
func (n *IpfsNode) HandlePeerFound(p pstore.PeerInfo) {
	log.Warning("trying peer info: ", p)
	ctx, cancel := context.WithTimeout(n.Context(), discoveryConnTimeout)
	defer cancel()
	if err := n.PeerHost.Connect(ctx, p); err != nil {
		log.Warning("Failed to connect to peer found by discovery: ", err)
	}
}

// startOnlineServicesWithHost  is the set of services which need to be
// initialized with the host and _before_ we start listening.
func (n *IpfsNode) startOnlineServicesWithHost(ctx context.Context, host p2phost.Host, routingOption RoutingOption, pubsub bool, ipnsps bool, iamOption IAMOption) error {
	// setup diagnostics service
	n.Ping = ping.NewPingService(host)

	if pubsub || ipnsps {
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}

		var service *floodsub.PubSub

		switch cfg.Pubsub.Router {
		case "":
			fallthrough
		case "floodsub":
			service, err = floodsub.NewFloodSub(ctx, host)

		case "gossipsub":
			service, err = floodsub.NewGossipSub(ctx, host)

		default:
			err = fmt.Errorf("Unknown pubsub router %s", cfg.Pubsub.Router)
		}

		if err != nil {
			return err
		}
		n.Floodsub = service
	}

	// setup routing service
	r, err := routingOption(ctx, host, n.Repo.Datastore(), n.RecordValidator)
	if err != nil {
		return err
	}
	n.Routing = r

	// TODO: I'm not a fan of type assertions like this but the
	// `RoutingOption` system doesn't currently provide access to the
	// IpfsNode.
	//
	// Ideally, we'd do something like:
	//
	// 1. Add some fancy method to introspect into tiered routers to extract
	//    things like the pubsub router or the DHT (complicated, messy,
	//    probably not worth it).
	// 2. Pass the IpfsNode into the RoutingOption (would also remove the
	//    PSRouter case below.
	// 3. Introduce some kind of service manager? (my personal favorite but
	//    that requires a fair amount of work).
	if dht, ok := r.(*dht.IpfsDHT); ok {
		n.DHT = dht
	}

	i, err := IamOption(ctx, host, n.Repo.Datastore())
	if err != nil {
		return err
	}

	n.IAM = i

	if ipnsps {
		n.PSRouter = psrouter.NewPubsubValueStore(
			ctx,
			host,
			n.Routing,
			n.Floodsub,
			n.RecordValidator,
		)
		n.Routing = rhelpers.Tiered{
			// Always check pubsub first.
			&rhelpers.Compose{
				ValueStore: &rhelpers.LimitedValueStore{
					ValueStore: n.PSRouter,
					Namespaces: []string{"ipns"},
				},
			},
			n.Routing,
		}
	}

	// Wrap standard peer host with routing system to allow unknown peer lookups
	n.PeerHost = rhost.Wrap(host, n.Routing)

	// set iam service
	//iam, err = iamOption(ctx, host, n.Repo.Datastore())
	//if err != nil {
	//	return err
	//}
	//n.IAM = iam

	// setup exchange service
	bitswapNetwork := bsnet.NewFromIpfsHost(n.PeerHost, n.Routing)
	n.Exchange = bitswap.New(ctx, bitswapNetwork, n.Blockstore)

	size, err := n.getCacheSize()
	if err != nil {
		return err
	}

	// setup name system
	n.Namesys = namesys.NewNameSystem(n.Routing, n.Repo.Datastore(), size)

	// setup ipns republishing
	return n.setupIpnsRepublisher()
}

// getCacheSize returns cache life and cache size
func (n *IpfsNode) getCacheSize() (int, error) {
	cfg, err := n.Repo.Config()
	if err != nil {
		return 0, err
	}

	cs := cfg.Ipns.ResolveCacheSize
	if cs == 0 {
		cs = 128
	}
	if cs < 0 {
		return 0, fmt.Errorf("cannot specify negative resolve cache size")
	}
	return cs, nil
}

func (n *IpfsNode) setupIpnsRepublisher() error {
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	n.IpnsRepub = ipnsrp.NewRepublisher(n.Namesys, n.Repo.Datastore(), n.PrivateKey, n.Repo.Keystore())

	if cfg.Ipns.RepublishPeriod != "" {
		d, err := time.ParseDuration(cfg.Ipns.RepublishPeriod)
		if err != nil {
			return fmt.Errorf("failure to parse config setting IPNS.RepublishPeriod: %s", err)
		}

		if !u.Debug && (d < time.Minute || d > (time.Hour*24)) {
			return fmt.Errorf("config setting IPNS.RepublishPeriod is not between 1min and 1day: %s", d)
		}

		n.IpnsRepub.Interval = d
	}

	if cfg.Ipns.RecordLifetime != "" {
		d, err := time.ParseDuration(cfg.Ipns.RepublishPeriod)
		if err != nil {
			return fmt.Errorf("failure to parse config setting IPNS.RecordLifetime: %s", err)
		}

		n.IpnsRepub.RecordLifetime = d
	}

	n.Process().Go(n.IpnsRepub.Run)

	return nil
}

// Process returns the Process object
func (n *IpfsNode) Process() goprocess.Process {
	return n.proc
}

// Close calls Close() on the Process object
func (n *IpfsNode) Close() error {
	return n.proc.Close()
}

// Context returns the IpfsNode context
func (n *IpfsNode) Context() context.Context {
	if n.ctx == nil {
		n.ctx = context.TODO()
	}
	return n.ctx
}

// teardown closes owned children. If any errors occur, this function returns
// the first error.
func (n *IpfsNode) teardown() error {
	log.Debug("core is shutting down...")
	// owned objects are closed in this teardown to ensure that they're closed
	// regardless of which constructor was used to add them to the node.
	var closers []io.Closer

	// NOTE: The order that objects are added(closed) matters, if an object
	// needs to use another during its shutdown/cleanup process, it should be
	// closed before that other object

	if n.FilesRoot != nil {
		closers = append(closers, n.FilesRoot)
	}

	if n.Exchange != nil {
		closers = append(closers, n.Exchange)
	}

	if n.Mounts.Ipfs != nil && !n.Mounts.Ipfs.IsActive() {
		closers = append(closers, mount.Closer(n.Mounts.Ipfs))
	}
	if n.Mounts.Ipns != nil && !n.Mounts.Ipns.IsActive() {
		closers = append(closers, mount.Closer(n.Mounts.Ipns))
	}

	if n.DHT != nil {
		closers = append(closers, n.DHT.Process())
	}

	if n.Blocks != nil {
		closers = append(closers, n.Blocks)
	}

	if n.Bootstrapper != nil {
		closers = append(closers, n.Bootstrapper)
	}

	if n.PeerHost != nil {
		closers = append(closers, n.PeerHost)
	}

	// Repo closed last, most things need to preserve state here
	closers = append(closers, n.Repo)

	var errs []error
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// OnlineMode returns whether or not the IpfsNode is in OnlineMode.
func (n *IpfsNode) OnlineMode() bool {
	return n.mode == onlineMode
}

// SetLocal will set the IpfsNode to local mode
func (n *IpfsNode) SetLocal(isLocal bool) {
	if isLocal {
		n.mode = localMode
	}
	n.localModeSet = true
}

// LocalMode returns whether or not the IpfsNode is in LocalMode
func (n *IpfsNode) LocalMode() bool {
	if !n.localModeSet {
		// programmer error should not happen
		panic("local mode not set")
	}
	return n.mode == localMode
}

// Bootstrap will set and call the IpfsNodes bootstrap function.
func (n *IpfsNode) Bootstrap(cfg BootstrapConfig) error {
	// TODO what should return value be when in offlineMode?
	if n.Routing == nil {
		return nil
	}

	if n.Bootstrapper != nil {
		n.Bootstrapper.Close() // stop previous bootstrap process.
	}

	// if the caller did not specify a bootstrap peer function, get the
	// freshest bootstrap peers from config. this responds to live changes.
	if cfg.BootstrapPeers == nil {
		cfg.BootstrapPeers = func() []pstore.PeerInfo {
			ps, err := n.loadBootstrapPeers()
			if err != nil {
				log.Warning("failed to parse bootstrap peers from config")
				return nil
			}
			return ps
		}
	}

	var err error
	n.Bootstrapper, err = Bootstrap(n, cfg)
	return err
}

func (n *IpfsNode) loadID() error {
	if n.Identity != "" {
		return errors.New("identity already loaded")
	}

	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	cid := cfg.Identity.PeerID
	if cid == "" {
		return errors.New("identity was not set in config (was 'ipfs init' run?)")
	}
	if len(cid) == 0 {
		return errors.New("no peer ID in config! (was 'ipfs init' run?)")
	}

	id, err := peer.IDB58Decode(cid)
	if err != nil {
		return fmt.Errorf("peer ID invalid: %s", err)
	}

	n.Identity = id
	return nil
}

// GetKey will return a key from the Keystore with name `name`.
func (n *IpfsNode) GetKey(name string) (ic.PrivKey, error) {
	if name == "self" {
		return n.PrivateKey, nil
	} else {
		return n.Repo.Keystore().Get(name)
	}
}

func (n *IpfsNode) LoadPrivateKey() error {
	if n.Identity == "" || n.Peerstore == nil {
		return errors.New("loaded private key out of order")
	}

	if n.PrivateKey != nil {
		return errors.New("private key already loaded")
	}

	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	sk, err := loadPrivateKey(&cfg.Identity, n.Identity)
	if err != nil {
		return err
	}

	n.PrivateKey = sk
	n.Peerstore.AddPrivKey(n.Identity, n.PrivateKey)
	n.Peerstore.AddPubKey(n.Identity, sk.GetPublic())
	return nil
}

func (n *IpfsNode) loadBootstrapPeers() ([]pstore.PeerInfo, error) {
	cfg, err := n.Repo.Config()
	if err != nil {
		return nil, err
	}

	parsed, err := cfg.BootstrapPeers()
	if err != nil {
		return nil, err
	}
	return toPeerInfos(parsed), nil
}

func (n *IpfsNode) loadFilesRoot() error {
	dsk := ds.NewKey("/local/filesroot")
	pf := func(ctx context.Context, c *cid.Cid) error {
		return n.Repo.Datastore().Put(dsk, c.Bytes())
	}

	var nd *merkledag.ProtoNode
	val, err := n.Repo.Datastore().Get(dsk)

	switch {
	case err == ds.ErrNotFound || val == nil:
		nd = ft.EmptyDirNode()
		err := n.DAG.Add(n.Context(), nd)
		if err != nil {
			return fmt.Errorf("failure writing to dagstore: %s", err)
		}
	case err == nil:
		c, err := cid.Cast(val)
		if err != nil {
			return err
		}

		rnd, err := n.DAG.Get(n.Context(), c)
		if err != nil {
			return fmt.Errorf("error loading filesroot from DAG: %s", err)
		}

		pbnd, ok := rnd.(*merkledag.ProtoNode)
		if !ok {
			return merkledag.ErrNotProtobuf
		}

		nd = pbnd
	default:
		return err
	}

	mr, err := mfs.NewRoot(n.Context(), n.DAG, nd, pf)
	if err != nil {
		return err
	}

	n.FilesRoot = mr
	return nil
}

// SetupOfflineRouting instantiates a routing system in offline mode. This is
// primarily used for offline ipns modifications.
func (n *IpfsNode) SetupOfflineRouting() error {
	if n.Routing != nil {
		// Routing was already set up
		return nil
	}

	// TODO: move this somewhere else.
	err := n.LoadPrivateKey()
	if err != nil {
		return err
	}

	n.Routing = offroute.NewOfflineRouter(n.Repo.Datastore(), n.RecordValidator)

	size, err := n.getCacheSize()
	if err != nil {
		return err
	}

	n.Namesys = namesys.NewNameSystem(n.Routing, n.Repo.Datastore(), size)

	return nil
}

func loadPrivateKey(cfg *config.Identity, id peer.ID) (ic.PrivKey, error) {
	sk, err := cfg.DecodePrivateKey("passphrase todo!")
	if err != nil {
		return nil, err
	}

	id2, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}

	if id2 != id {
		return nil, fmt.Errorf("private key in config does not match id: %s != %s", id, id2)
	}

	return sk, nil
}

func listenAddresses(cfg *config.Config) ([]ma.Multiaddr, error) {
	var listen []ma.Multiaddr
	for _, addr := range cfg.Addresses.Swarm {
		maddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("failure to parse config.Addresses.Swarm: %s", cfg.Addresses.Swarm)
		}
		listen = append(listen, maddr)
	}

	return listen, nil
}

type ConstructPeerHostOpts struct {
	AddrsFactory      p2pbhost.AddrsFactory
	DisableNatPortMap bool
	DisableRelay      bool
	EnableRelayHop    bool
	ConnectionManager ifconnmgr.ConnManager
}

type HostOption func(ctx context.Context, id peer.ID, ps pstore.Peerstore, options ...libp2p.Option) (p2phost.Host, error)

var DefaultHostOption HostOption = constructPeerHost

// isolates the complex initialization steps
func constructPeerHost(ctx context.Context, id peer.ID, ps pstore.Peerstore, options ...libp2p.Option) (p2phost.Host, error) {
	pkey := ps.PrivKey(id)
	if pkey == nil {
		return nil, fmt.Errorf("missing private key for node ID: %s", id.Pretty())
	}
	options = append([]libp2p.Option{libp2p.Identity(pkey), libp2p.Peerstore(ps)}, options...)
	return libp2p.New(ctx, options...)
}

func filterRelayAddrs(addrs []ma.Multiaddr) []ma.Multiaddr {
	var raddrs []ma.Multiaddr
	for _, addr := range addrs {
		_, err := addr.ValueForProtocol(circuit.P_CIRCUIT)
		if err == nil {
			continue
		}
		raddrs = append(raddrs, addr)
	}
	return raddrs
}

func composeAddrsFactory(f, g p2pbhost.AddrsFactory) p2pbhost.AddrsFactory {
	return func(addrs []ma.Multiaddr) []ma.Multiaddr {
		return f(g(addrs))
	}
}

// startListening on the network addresses
func startListening(host p2phost.Host, cfg *config.Config) error {
	listenAddrs, err := listenAddresses(cfg)
	if err != nil {
		return err
	}

	// Actually start listening:
	if err := host.Network().Listen(listenAddrs...); err != nil {
		return err
	}

	// list out our addresses
	addrs, err := host.Network().InterfaceListenAddresses()
	if err != nil {
		return err
	}
	log.Infof("Swarm listening at: %s", addrs)
	return nil
}

func constructDHTRouting(ctx context.Context, host p2phost.Host, dstore ds.Batching, validator record.Validator) (routing.IpfsRouting, error) {
	return dht.New(
		ctx, host,
		dhtopts.Datastore(dstore),
		dhtopts.Validator(validator),
	)
}

func constructClientDHTRouting(ctx context.Context, host p2phost.Host, dstore ds.Batching, validator record.Validator) (routing.IpfsRouting, error) {
	return dht.New(
		ctx, host,
		dhtopts.Client(true),
		dhtopts.Datastore(dstore),
		dhtopts.Validator(validator),
	)
}

func constructIAM(ctx context.Context, host p2phost.Host, ds ds.Batching) (*iam.IAM, error) {
	fmt.Println("construct IAM")
	return iam.NewIAMService(ctx, host, ds), nil
}

type RoutingOption func(context.Context, p2phost.Host, ds.Batching, record.Validator) (routing.IpfsRouting, error)

type DiscoveryOption func(context.Context, p2phost.Host) (discovery.Service, error)

type IAMOption func(context.Context, p2phost.Host, ds.Batching) (*iam.IAM, error)

var DHTOption RoutingOption = constructDHTRouting
var DHTClientOption RoutingOption = constructClientDHTRouting
var NilRouterOption RoutingOption = nilrouting.ConstructNilRouting

var IamOption IAMOption = constructIAM
