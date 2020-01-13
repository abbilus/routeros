package routeros

import (
	"bitosismessages/log"
	"fmt"

	errors "github.com/pkg/errors"
)

type Ros struct {
	id        int
	bitosisID int
	username  string
	password  string
	address   string
	port      int
	client    *Client
	Verbose   bool
}

func (ros *Ros) SetID(id int) {
	ros.id = id
}

func (ros *Ros) GetID() int {
	return ros.id
}

func (ros *Ros) SetBitosisID(id int) {
	ros.bitosisID = id
}

func (ros *Ros) GetBitosisID() int {
	return ros.bitosisID
}

func (ros *Ros) SetUsername(username string) {
	ros.username = username
}

func (ros *Ros) SetPassword(password string) {
	ros.password = password
}

func (ros *Ros) SetAddress(address string) {
	ros.address = address
}

func (ros *Ros) GetAddress() string {
	return ros.address
}

func (ros *Ros) SetPort(port int) {
	ros.port = port
}

func (ros *Ros) Connect(username, password, address string, port int) (err error) {
	ros.username = username
	ros.password = password
	ros.address = address
	ros.port = port
	err = ros.connect()
	return
}

func (ros *Ros) connect() (err error) {
	addr := ros.address
	if ros.port != 0 {
		addr = fmt.Sprintf("%s:%d", ros.address, ros.port)
	}
	log.Debugf("%s connecting", addr)
	ros.client, err = Dial(addr, ros.username, ros.password)
	// ros.client, err = routeros.Dial(addr, ros.username, ros.password)
	if err != nil {
		errors.Wrapf(err, "%s: ros/connect", ros.address)
		errors.WithStack(err)
		return
	}
	log.Debugf("Host %s connected", ros.GetAddress())
	return
}

func (ros *Ros) Reconnect() (err error) {
	ros.connect()
	return
}

func (ros *Ros) AddListEntry(list string, ip string) (err error) {
	if list == "" {
		err = fmt.Errorf("Ros/AddListEntry: list can't be empty")
		return
	}
	if ip == "" {
		err = fmt.Errorf("Ros/AddListEntry: ip can't be empty")
		return
	}
	cmd := "/ip/firewall/address-list/add"
	log.Debugf("Preparing to add %s to list %s\n", ip, list)
	var r *Reply
	if ros.client != nil {
		r, err = ros.client.Run(cmd, fmt.Sprintf("=list=%s", list), fmt.Sprintf("=address=%s", ip))
	} else {
		err = fmt.Errorf("Ros client %s not connected", ros.GetAddress())
	}
	log.Dump(r)
	if err != nil {
		errors.Wrapf(err, "%s: ros/GetSnapshot", ros.address)
		errors.WithStack(err)
		return
	}
	return
}

func (ros *Ros) DelListEntries(list string) (err error) {
	if list == "" {
		err = fmt.Errorf("Ros/DelListEntries: list can't be empty")
		return
	}
	cmd := "/ip/firewall/address-list/print"
	var r *Reply
	if ros.client != nil {
		r, err = ros.client.Run(cmd, fmt.Sprintf(`?list=%s`, list))
		for _, re := range r.Re {
			var id string
			cmd := "/ip/firewall/address-list/remove"
			id = re.Map[".id"]
			_, err = ros.client.Run(cmd, fmt.Sprintf(`=.id=%s`, id))
			if err != nil {
				log.Error(err)
			}
		}
	} else {
		err = fmt.Errorf("Ros client %s not connected", ros.GetAddress())
	}
	if err != nil {
		errors.Wrapf(err, "%s: ros/GetSnapshot", ros.address)
		errors.WithStack(err)
		return
	}
	return
}

func (ros *Ros) GetFirewallListsByName(listName string) (err error) {
	cmd := "/ip/firewall/address-list/print"
	var r *Reply
	if ros.client != nil {
		r, err = ros.client.Run(cmd, fmt.Sprintf(`?list=%s`, listName))
		for _, re := range r.Re {
			fmt.Printf("%+v", re)
		}
	} else {
		err = fmt.Errorf("Ros client %s not connected", ros.GetAddress())
	}
	if err != nil {
		errors.WithStack(err)
		return
	}
	return
}

func (ros *Ros) DelListEntry(list string, ip string) (err error) {
	if list == "" {
		err = fmt.Errorf("Ros/DelListEntry: list can't be empty")
		return
	}
	if ip == "" {
		err = fmt.Errorf("Ros/DelListEntry: ip can't be empty")
		return
	}
	cmd := "/ip/firewall/address-list/print"
	var r *Reply
	if ros.client != nil {
		r, err = ros.client.Run(cmd, fmt.Sprintf(`?list=%s`, list), fmt.Sprintf(`?address=%s`, ip))
		for _, re := range r.Re {
			var id string
			cmd := "/ip/firewall/address-list/remove"
			id = re.Map[".id"]
			_, err = ros.client.Run(cmd, fmt.Sprintf(`=.id=%s`, id))
			if err != nil {
				log.Error(err)
			}
		}
	} else {
		err = fmt.Errorf("Ros client %s not connected", ros.GetAddress())
	}
	if err != nil {
		errors.Wrapf(err, "%s: ros/GetSnapshot", ros.address)
		errors.WithStack(err)
		return
	}
	return
}

func (ros *Ros) GetMacByIP(ipaddr string) (mac string, err error) {
	cmd := "/ip/arp/print"
	var r *Reply
	if ros.client != nil {
		r, err = ros.client.Run(cmd, fmt.Sprintf(`?address=%s`, ipaddr))
		if len(r.Re) == 1 {
			mac = r.Re[0].Map["mac-address"]
		}
	} else {
		err = fmt.Errorf("Ros client %s not connected", ros.GetAddress())
	}
	if err != nil {
		errors.Wrapf(err, "%s: ros/GetMacByIP", ros.address)
		errors.WithStack(err)
		return
	}
	return
}

func (ros *Ros) Disconnect() {
	ros.client.Close()
	log.Debugf("Host %s disconnected", ros.address)
}
