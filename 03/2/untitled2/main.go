// main.go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:sa123456@tcp(127.0.0.1:3306)/ry?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接失败: " + err.Error())
	}

	db.AutoMigrate(&Account{}, &Transaction{})

	db.Exec("TRUNCATE TABLE transactions")
	db.Exec("TRUNCATE TABLE accounts")
	db.Create(&Account{Balance: 500.0})
	db.Create(&Account{Balance: 200.0})

	fromID := uint(1)
	toID := uint(2)
	amount := 100.0

	err = transfer(db, fromID, toID, amount)
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	} else {
		fmt.Println("转账成功！")
	}

	printAccounts(db)
	printTransactions(db)
}

func transfer(db *gorm.DB, fromID, toID uint, amount float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var fromAcc Account
		if err := tx.First(&fromAcc, fromID).Error; err != nil {
			return fmt.Errorf("转出账户不存在")
		}

		if fromAcc.Balance < amount {
			return fmt.Errorf("余额不足：当前 %.2f，需 %.2f", fromAcc.Balance, amount)
		}

		tx.Model(&fromAcc).Update("balance", gorm.Expr("balance - ?", amount))

		var toAcc Account
		tx.First(&toAcc, toID)
		tx.Model(&toAcc).Update("balance", gorm.Expr("balance + ?", amount))

		transaction := Transaction{
			FromAccountID: fromID,
			ToAccountID:   toID,
			Amount:        amount,
		}
		tx.Create(&transaction)

		return nil
	})
}

func printAccounts(db *gorm.DB) {
	var accounts []Account
	db.Find(&accounts)
	fmt.Println("\n=== 账户余额 ===")
	for _, a := range accounts {
		fmt.Printf("ID: %d, 余额: %.2f\n", a.ID, a.Balance)
	}
}

func printTransactions(db *gorm.DB) {
	var txs []Transaction
	db.Find(&txs)
	fmt.Println("\n=== 交易记录 ===")
	for _, t := range txs {
		fmt.Printf("转账: %d → %d, 金额: %.2f\n", t.FromAccountID, t.ToAccountID, t.Amount)
	}
}
