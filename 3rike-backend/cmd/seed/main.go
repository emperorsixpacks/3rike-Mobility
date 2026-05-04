// cmd/seed/main.go — seeds demo data for 3riKE.
// Run: DATABASE_URL=... go run ./cmd/seed
package main

import (
	"log"
	"os"

	"github.com/3rike12/3rike-backend/internal/repository"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	if err := repository.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// --- Users ---
	users := []repository.User{
		{Email: "emeka@3rike.xyz", PasswordHash: "$2a$10$seed", Role: "driver"},
		{Email: "kwame@3rike.xyz", PasswordHash: "$2a$10$seed", Role: "driver"},
		{Email: "jean@3rike.xyz", PasswordHash: "$2a$10$seed", Role: "driver"},
		{Email: "amara@3rike.xyz", PasswordHash: "$2a$10$seed", Role: "investor"},
		{Email: "fatima@3rike.xyz", PasswordHash: "$2a$10$seed", Role: "investor"},
	}
	for i := range users {
		db.Where(repository.User{Email: users[i].Email}).FirstOrCreate(&users[i])
		log.Printf("✅ user: %s (id=%d)", users[i].Email, users[i].ID)
	}

	// --- Drivers ---
	drivers := []repository.Driver{
		{UserID: users[0].ID, FullName: "Emeka Okafor", Phone: "+2348012345678", Country: "Nigeria", CreditScore: 720, WeeksRemaining: 58},
		{UserID: users[1].ID, FullName: "Kwame Asante", Phone: "+233244123456", Country: "Ghana", CreditScore: 650, WeeksRemaining: 42},
		{UserID: users[2].ID, FullName: "Jean-Pierre Habimana", Phone: "+250788123456", Country: "Rwanda", CreditScore: 580, WeeksRemaining: 65},
	}
	for i := range drivers {
		db.Where(repository.Driver{UserID: drivers[i].UserID}).FirstOrCreate(&drivers[i])
		log.Printf("✅ driver: %s (id=%d)", drivers[i].FullName, drivers[i].ID)
	}

	// --- Tricycles ---
	tricycles := []repository.Tricycle{
		{DriverID: drivers[0].ID, Make: "Bajaj", Model: "RE 4S", IsEV: false, PriceUSD: 1800, Status: "financing", TotalFractions: 0},
		{DriverID: drivers[1].ID, Make: "TVS", Model: "King Deluxe", IsEV: false, PriceUSD: 1650, Status: "tokenized", ContractID: "stub-contract-tvs-king"},
		{DriverID: drivers[2].ID, Make: "Piaggio", Model: "Ape City", IsEV: false, PriceUSD: 2100, Status: "fractionalized", ContractID: "stub-contract-piaggio-ape", TotalFractions: 100},
		{DriverID: drivers[0].ID, Make: "Keke EV", Model: "E-Trike Pro", IsEV: true, PriceUSD: 2800, Status: "available"},
		{DriverID: drivers[1].ID, Make: "Bajaj", Model: "RE Compact", IsEV: false, PriceUSD: 1500, Status: "financing"},
	}
	for i := range tricycles {
		db.Where(repository.Tricycle{DriverID: tricycles[i].DriverID, Model: tricycles[i].Model}).FirstOrCreate(&tricycles[i])
		log.Printf("✅ tricycle: %s %s (id=%d, status=%s)", tricycles[i].Make, tricycles[i].Model, tricycles[i].ID, tricycles[i].Status)
	}

	// --- Investors ---
	investors := []repository.Investor{
		{UserID: users[3].ID, FullName: "Amara Nwosu", WalletAddress: "0xabc123"},
		{UserID: users[4].ID, FullName: "Fatima Al-Hassan", WalletAddress: "0xdef456"},
	}
	for i := range investors {
		db.Where(repository.Investor{UserID: investors[i].UserID}).FirstOrCreate(&investors[i])
		log.Printf("✅ investor: %s (id=%d)", investors[i].FullName, investors[i].ID)
	}

	// --- Fractions (investments in the fractionalized tricycle) ---
	fractions := []repository.Fraction{
		{TricycleID: tricycles[2].ID, InvestorID: investors[0].ID, Units: 30, PricePerUnit: 21.00},
		{TricycleID: tricycles[2].ID, InvestorID: investors[1].ID, Units: 20, PricePerUnit: 21.00},
	}
	for i := range fractions {
		db.Where(repository.Fraction{TricycleID: fractions[i].TricycleID, InvestorID: fractions[i].InvestorID}).FirstOrCreate(&fractions[i])
		log.Printf("✅ fraction: investor=%d tricycle=%d units=%d", fractions[i].InvestorID, fractions[i].TricycleID, fractions[i].Units)
	}

	// --- Payments (3 weeks of payments for driver 1) ---
	for week := 1; week <= 3; week++ {
		p := repository.Payment{
			DriverID: drivers[0].ID, TricycleID: tricycles[0].ID,
			AmountLocal: 15000, AmountUSDC: 10.50, Currency: "NGN",
			Status: "confirmed", WeekNumber: week,
		}
		db.Where(repository.Payment{DriverID: p.DriverID, WeekNumber: week}).FirstOrCreate(&p)
		log.Printf("✅ payment: driver=%d week=%d", p.DriverID, week)
	}

	// --- Savings ---
	savings := repository.Savings{DriverID: drivers[0].ID, BalanceUSDC: 45.00}
	db.Where(repository.Savings{DriverID: savings.DriverID}).FirstOrCreate(&savings)
	log.Printf("✅ savings: driver=%d balance=%.2f USDC", savings.DriverID, savings.BalanceUSDC)

	log.Println("🎉 seed complete")
}
