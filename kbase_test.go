package kbase
import (
	"testing"
	"time"
	"github.com/gl-works/kbase"
	"strconv"
	"math/rand"
	"fmt"
)

const (
	clustersize = 10
	testportstart = 9000
	riskRandomness = false
	space = "test"
)

var (
	keys = []string{"key1", "key2", "key3", "key4", "key5"}
	instances []*kbase.Kbase
)

func serialize(value kbase.Value) ([]byte, error) {
	return []byte(value.(string)), nil
}

func deserialize(bytes []byte) (kbase.Value, error)  {
	return string(bytes), nil
}

func init() {
	hosts := []string{}
	for i:=0; i<clustersize; i++ {
		hosts = append(hosts, "localhost:"+strconv.Itoa(testportstart + i*2))
	}

	instances = []*kbase.Kbase{}
	for i, host := range hosts {
		instances = append(instances, kbase.NewKbase(host, hosts[0:i+1],
			&kbase.KspaceDef{Name:space, Serializer:serialize, Deserializer:deserialize}))
	}
	for _, instance := range instances {
		for instance.ClusterSize() < clustersize {
			time.Sleep(time.Millisecond * 1)
		}
		if riskRandomness {
			break
		}
	}
	fmt.Println("Done init test cluster")
}

func Test_kbase(t *testing.T) {
	t.Logf("Starting test sub")
	for _,key := range keys {
		message := "any random message " + strconv.Itoa(rand.Int())
		if err := instances[rand.Int() % clustersize].SpaceOf(space).Offer(key, message); err != nil {
			t.Errorf("Offer error: %s", err)
			t.Fail()
		} else {
			for _, anyInstance := range instances {
				if values := anyInstance.SpaceOf(space).Poll(key); len(values) != 1 {
					t.Errorf("unexpected 1 value got %d", len(values))
					t.Fail()
				} else {
					for _, value := range values {
						if value != message {
							t.Errorf("expect %v for key %s got %v",
								message, key, value)
							t.Fail()
						}
					}
				}
			}
		}
	}
	t.Log("OK")
}
