package tests

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"shop_srvs/inventory_srv/proto"
	"sync"
	"testing"
)

var invClient proto.InventoryClient

func Init() {
	var err error
	conn, err := grpc.Dial("192.168.112.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	invClient = proto.NewInventoryClient(conn)
}

func TestSetInv(t *testing.T) {
	Init()
	resp, err := invClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 422,
		Num:     101,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestInvDetail(t *testing.T) {
	Init()
	resp, err := invClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 421,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestSell(t *testing.T) {
	Init()
	resp, err := invClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInvInfos: []*proto.GoodsInvInfo{
			{
				GoodsId: 421,
				Num:     15,
			},
			{
				GoodsId: 422,
				Num:     15,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestReback(t *testing.T) {
	Init()
	resp, err := invClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInvInfos: []*proto.GoodsInvInfo{
			{
				GoodsId: 421,
				Num:     1,
			},
			{
				GoodsId: 422,
				Num:     12,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestAddInv(t *testing.T) {
	Init()
	var i int
	for i = 421; i <= 840; i++ {
		resp, err := invClient.SetInv(context.Background(), &proto.GoodsInvInfo{
			GoodsId: int32(i),
			Num:     100,
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(resp)
	}
}

func SellCur(wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := invClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInvInfos: []*proto.GoodsInvInfo{
			{
				GoodsId: 421,
				Num:     1,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.String())
}

func TestConcurrency(t *testing.T) {
	Init()
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go SellCur(&wg)
	}
	wg.Wait()
}
