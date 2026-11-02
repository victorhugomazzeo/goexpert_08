package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateAuction_ClosesAutomatically(t *testing.T) {
	t.Setenv("AUCTION_INTERVAL", "1s")

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/?authSource=admin"))
	if err != nil {
		t.Fatalf("erro ao conectar: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("mongo não está acessível: %v", err)
	}

	db := client.Database("auctions_test")
	defer db.Drop(ctx)

	repo := NewAuctionRepository(db)

	auctionEntity, internalErr := auction_entity.CreateAuction("Notebook", "Eletronicos", "Notebook usado em bom estado", auction_entity.Used)
	if internalErr != nil {
		t.Fatalf("erro ao criar entidade: %v", internalErr)
	}
	if internalErr := repo.CreateAuction(ctx, auctionEntity); internalErr != nil {
		t.Fatalf("erro ao salvar leilão: %v", internalErr)
	}

	found, internalErr := repo.FindAuctionById(ctx, auctionEntity.Id)
	if internalErr != nil {
		t.Fatalf("erro ao buscar leilão: %v", internalErr)
	}
	if found.Status != auction_entity.Active {
		t.Fatalf("esperado Active, recebido %v", found.Status)
	}

	time.Sleep(2 * time.Second)

	found, internalErr = repo.FindAuctionById(ctx, auctionEntity.Id)
	if internalErr != nil {
		t.Fatalf("erro ao buscar leilão: %v", internalErr)
	}
	if found.Status != auction_entity.Completed {
		t.Fatalf("esperado Completed, recebido %v", found.Status)
	}

}
