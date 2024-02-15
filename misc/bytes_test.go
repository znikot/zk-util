package misc

import (
	"fmt"
	"testing"
)

func TestIntToBytes(t *testing.T) {
	fmt.Printf("%x \n", IntToBytes(1))
	fmt.Printf("%x \n", IntToBytes(-1))
}

func TestBytesToInt(t *testing.T) {
	fmt.Printf("%d \n", BytesToInt([]byte{0, 0, 0, 0x01}))
	fmt.Printf("%d \n", BytesToInt([]byte{0xff, 0xff, 0xff, 0xff}))
}

func TestInt64(t *testing.T) {
	fmt.Printf("%X \n", Int64ToBytes(1))
	fmt.Printf("%X \n", Int64ToBytes(-1))
	fmt.Printf("%X \n", Int64ToBytes(1987))
	fmt.Printf("%X \n", Int64ToBytes(-1987))

	fmt.Println("=====================")
	var bytes []byte
	bytes = Int64ToBytes(1)
	fmt.Printf("%d\t->\t%X \n", BytesToInt64(bytes), Int64ToBytes(BytesToInt64(bytes)))
	bytes = Int64ToBytes(-1)
	fmt.Printf("%d\t->\t%X \n", BytesToInt64(bytes), Int64ToBytes(BytesToInt64(bytes)))
	bytes = Int64ToBytes(1987)
	fmt.Printf("%d\t->\t%X \n", BytesToInt64(bytes), Int64ToBytes(BytesToInt64(bytes)))
	bytes = Int64ToBytes(-1987)
	fmt.Printf("%d\t->\t%X \n", BytesToInt64(bytes), Int64ToBytes(BytesToInt64(bytes)))
}

func TestInt(t *testing.T) {
	fmt.Printf("%X \n", IntToBytes(1))
	fmt.Printf("%X %X\n", IntToBytes(-1), -18761)
	fmt.Printf("%X \n", IntToBytes(1987))
	fmt.Printf("%X \n", IntToBytes(-1987))
	fmt.Println("=====================")
	var bytes []byte
	bytes = IntToBytes(1)
	fmt.Printf("%d\t->\t%X \n", BytesToInt(bytes), IntToBytes(BytesToInt(bytes)))
	bytes = IntToBytes(-1)
	fmt.Printf("%d\t->\t%X \n", BytesToInt(bytes), IntToBytes(BytesToInt(bytes)))
	bytes = IntToBytes(1987)
	fmt.Printf("%d\t->\t%X \n", BytesToInt(bytes), IntToBytes(BytesToInt(bytes)))
	bytes = IntToBytes(-1987)
	fmt.Printf("%d\t->\t%X \n", BytesToInt(bytes), IntToBytes(BytesToInt(bytes)))
}
