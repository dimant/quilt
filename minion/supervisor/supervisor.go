package supervisor

import (
	"fmt"

	"github.com/NetSys/di/db"
	"github.com/NetSys/di/minion/docker"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("supervisor")

type supervisor struct {
	dk docker.Client

	role      db.Role
	etcdToken string
	leaderIP  string
	IP        string
	leader    bool
}

func Run(conn db.Conn, dk docker.Client) {
	sv := supervisor{dk: dk, role: db.None}
	for range conn.Trigger(db.MinionTable).C {
		var minion db.Minion
		minions := conn.SelectFromMinion(nil)
		if len(minions) == 1 {
			minion = minions[0]
		}

		if sv.role == minion.Role &&
			sv.etcdToken == minion.EtcdToken &&
			sv.leaderIP == minion.LeaderIP &&
			sv.IP == minion.PrivateIP &&
			sv.leader == minion.Leader {
			continue
		}

		if minion.Role != sv.role {
			sv.dk.RemoveAll()
		}

		switch minion.Role {
		case db.Master:
			sv.updateMaster(minion.PrivateIP, minion.EtcdToken,
				minion.Leader)
		case db.Worker:
			sv.updateWorker(minion.PrivateIP, minion.LeaderIP,
				minion.EtcdToken)
		}

		sv.role = minion.Role
		sv.etcdToken = minion.EtcdToken
		sv.leaderIP = minion.LeaderIP
		sv.IP = minion.PrivateIP
		sv.leader = minion.Leader
	}
}

func (sv supervisor) updateWorker(IP, leaderIP, etcdToken string) {
	if sv.etcdToken != etcdToken {
		sv.dk.Remove(docker.Etcd)
	}

	if sv.leaderIP != leaderIP || sv.IP != IP {
		sv.dk.Remove(docker.Kubelet)
	}

	etcdArgs := []string{"--discovery=" + etcdToken, "--proxy=on"}
	sv.dk.Run(docker.Etcd, etcdArgs)
	sv.dk.Run(docker.Ovsdb, nil)
	sv.dk.Run(docker.Ovncontroller, nil)
	sv.dk.Run(docker.Ovsvswitchd, nil)

	if leaderIP == "" || IP == "" {
		return
	}

	args := []string{"/usr/bin/boot-worker", IP, leaderIP}
	sv.dk.Run(docker.Kubelet, args)

	sv.dk.Exec(docker.Ovsvswitchd, []string{"ovs-vsctl", "set", "Open_vSwitch", ".",
		fmt.Sprintf("external_ids:ovn-remote=\"tcp:%s:6640\"", leaderIP),
		fmt.Sprintf("external_ids:ovn-encap-ip=%s", IP),
		"external_ids:ovn-encap-type=\"geneve\"",
	})
}

func (sv supervisor) updateMaster(IP, etcdToken string, leader bool) {
	if sv.IP != IP || sv.etcdToken != etcdToken {
		sv.dk.Remove(docker.Etcd)
	}

	if sv.IP != IP {
		sv.dk.Remove(docker.Kubelet)
	}

	if IP == "" || etcdToken == "" {
		return
	}

	etcdArgs := []string{
		fmt.Sprintf("--name=master-%s", IP),
		fmt.Sprintf("--discovery=%s", etcdToken),
		fmt.Sprintf("--advertise-client-urls=http://%s:2379", IP),
		fmt.Sprintf("--listen-peer-urls=http://%s:2380", IP),
		fmt.Sprintf("--initial-advertise-peer-urls=http://%s:2380", IP),
		"--listen-client-urls=http://0.0.0.0:2379",
	}
	sv.dk.Run(docker.Etcd, etcdArgs)
	sv.dk.Run(docker.Ovsdb, nil)
	sv.dk.Run(docker.Kubelet, []string{"/usr/bin/boot-master", IP})

	if leader {
		/* XXX: If we fail to boot ovn-northd, we should give up
		* our leadership somehow.  This ties into the general
		* problem of monitoring health. */
		sv.dk.Run(docker.Ovnnorthd, nil)
	} else {
		sv.dk.Remove(docker.Ovnnorthd)
	}
}