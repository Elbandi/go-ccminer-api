package ccminer

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"errors"
	"github.com/mitchellh/mapstructure"
)

type CCMiner struct {
	server string
}

// New returns a CGMiner pointer, which is used to communicate with a running
// CGMiner instance. Note that New does not attempt to connect to the miner.
func New(hostname string, port int64) *CCMiner {
	miner := new(CCMiner)
	server := fmt.Sprintf("%s:%d", hostname, port)
	miner.server = server

	return miner
}

func (miner *CCMiner) runCommand(command, argument string) (string, error) {
	conn, err := net.Dial("tcp", miner.server)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s|%s", command, argument)
	result, err := bufio.NewReader(conn).ReadString('\x00')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(result, "\x00"), nil
}

type Summary struct {
	Name            string `mapstructure:"NAME"`
	Version         string `mapstructure:"VER"`
	Api             string `mapstructure:"API"`
	Algo            string `mapstructure:"ALGO"`
	Gpus            uint8 `mapstructure:"GPUS"`
	Khs             float32 `mapstructure:"KHS"`
	Solved          float32 `mapstructure:"SOLV"`
	Accepted        uint64 `mapstructure:"ACC"`
	Rejected        uint64 `mapstructure:"REJ"`
	AcceptedMinutes float32 `mapstructure:"ACCMN"`
	Diff            float32 `mapstructure:"DIFF"`
	NetKhs          float32 `mapstructure:"NETKHS"`
	Pools           uint8 `mapstructure:"POOLS"`
	Wait            uint32 `mapstructure:"WAIT"`
	Uptime          uint32 `mapstructure:"UPTIME"`
	TimeStamp       uint64 `mapstructure:"UPTIME"`
}

type Device struct {
	Id             uint8 `mapstructure:"GPU"`
	Bus            uint8 `mapstructure:"BUS"`
	Card           string `mapstructure:"CARD"`
	Temp           float32 `mapstructure:"TEMP"`
	Power          uint32 `mapstructure:"POWER"`
	Fan            uint16 `mapstructure:"FAN"`
	Rpm            uint16 `mapstructure:"RPM"`
	GpuFreq        uint16 `mapstructure:"FREQ"`
	MemFreq        uint16 `mapstructure:"MEMFREQ"`
	MonitorGpuFreq uint16 `mapstructure:"GPUF"`
	MonitorMemFreq uint16 `mapstructure:"MEMF"`
	Khs            float32 `mapstructure:"KHS"`
	KhsWatt        float32 `mapstructure:"KHW"`
	PowerLimit     uint32 `mapstructure:"PLIM"`
	Accepted       uint64 `mapstructure:"ACC"`
	Rejected       uint64 `mapstructure:"REJ"`
	HardwareErrors uint64 `mapstructure:"HWF"`
	Intensity      string `mapstructure:"I"`
	Throughput     uint64 `mapstructure:"THR"`
}

type Pool struct {
	Name          string `mapstructure:"POOL"`
	Algo          string `mapstructure:"ALGO"`
	Url           string `mapstructure:"URL"`
	User          string `mapstructure:"USER"`
	Solved        float32 `mapstructure:"SOLV"`
	Accepted      uint64 `mapstructure:"ACC"`
	Rejected      uint64 `mapstructure:"REJ"`
	Stale         uint64 `mapstructure:"STALE"`
	Height        uint64 `mapstructure:"H"`
	LastJobId     string `mapstructure:"JOB"`
	Diff          float32 `mapstructure:"DIFF"`
	BestShareDiff float32 `mapstructure:"BEST"`
	XNonce2Size   uint8 `mapstructure:"N2SZ"`
	Nonce2        string `mapstructure:"N2"`
	Ping          float32 `mapstructure:"PING"`
	Disconnects   uint16 `mapstructure:"DISCO"`
	Wait          uint32 `mapstructure:"WAIT"`
	Uptime        uint32 `mapstructure:"UPTIME"`
	TimeStamp     uint64 `mapstructure:"LAST"`
}

func convertToStruct(data string, rawVal interface{}) error {
	ss := strings.Split(data, ";")
	if len(ss) == 0 {
		return errors.New("Empty data")
	}
	m := make(map[string]string)
	for _, pair := range ss {
		z := strings.Split(pair, "=")
		m[z[0]] = z[1]
	}
	err := mapstructure.WeakDecode(m, rawVal)
	if err != nil {
		return err
	}
	return nil
}

// Summary returns basic information on the miner. See the Summary struct.
func (miner *CCMiner) Summary() (*Summary, error) {
	result, err := miner.runCommand("summary", "")
	if err != nil {
		return nil, err
	}
	s := strings.Split(result, "|")
	if len(s) != 2 {
		return nil, errors.New("Not valid ccminer summary response")
	}

	var summary Summary
	err = convertToStruct(result, &summary)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

// Devs returns basic information on the miner. See the Devs struct.
func (miner *CCMiner) Devs() ([]Device, error) {
	result, err := miner.runCommand("threads", "")
	if err != nil {
		return nil, err
	}
	s := strings.Split(result, "|")
	slen := len(s)
	if slen < 2 {
		return nil, errors.New("Not valid ccminer threads response")
	}
	var devs[] Device
	for _, val := range s[0:slen - 1] {
		var dev Device
		err = convertToStruct(val, &dev)
		if err != nil {
			return nil, err
		}
		devs = append(devs, dev)
	}
	return devs, err
}

// Pools returns a slice of Pool structs, one per pool.
func (miner *CCMiner) Pools() ([]Pool, error) {
	result, err := miner.runCommand("pool", "")
	if err != nil {
		return nil, err
	}
	fmt.Print(result)
	s := strings.Split(result, "|")
	slen := len(s)
	if slen < 2 {
		return nil, errors.New("Not valid ccminer pool response")
	}
	var pools[] Pool
	for _, val := range s[0:slen - 1] {
		var pool Pool
		err = convertToStruct(val, &pool)
		if err != nil {
			return nil, err
		}
		pools = append(pools, pool)
	}
	return pools, err
}