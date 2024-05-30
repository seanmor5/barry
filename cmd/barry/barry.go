package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/seanmor5/barry/internal/mercury"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	apiKey            string
	filterAccountsStr string
	startDate         string
	endDate           string
	aggregate         string
	counterparty      string
)

var rootCmd = &cobra.Command{
	Use:   "barry",
	Short: "Barry is a CLI for performing common accounting and banking tasks",
}

var balancesCmd = &cobra.Command{
	Use:   "balances",
	Short: "View account balances",
	Run: func(cmd *cobra.Command, args []string) {
		if apiKey == "" {
			apiKey = viper.GetString("MERCURY_API_KEY")
		}

		if apiKey == "" {
			log.Fatalf("You must provide an API key, either by setting MERCURY_API_KEY or passing --api-key")
		}

		var filterAccounts []string
		if filterAccountsStr != "" {
			filterAccounts = strings.Split(filterAccountsStr, ",")
		}

		config := mercury.Config{
			APIKey: apiKey,
		}

		accountsResponse, err := mercury.ListAccounts(config)
		if err != nil {
			log.Fatalf("Something went wrong listing accounts: %v", err)
		}

		for _, account := range accountsResponse.Accounts {
			if filterAccountsStr == "" || contains(filterAccounts, account.ID) {
				fmt.Printf("%s (%s)\nCurrent Balance: %2f, Available Balance: %2f\n", account.Name, account.ID, account.CurrentBalance, account.AvailableBalance)
			}
		}
	},
}

var spendCmd = &cobra.Command{
	Use:   "spend",
	Short: "Track spend across counterparties and periods",
	Run: func(cmd *cobra.Command, args []string) {
		if apiKey == "" {
			apiKey = viper.GetString("MERCURY_API_KEY")
		}

		if apiKey == "" {
			log.Fatalf("You must provide an API key, either by setting MERCURY_API_KEY or passing --api-key")
		}

		listTransactionParams := &mercury.ListTransactionsParams{}

		if startDate != "" {
			_, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				log.Fatalf("Invalid start date format. Please use YYYY-MM-DD.")
			}
			listTransactionParams.Start = &startDate
		}

		if endDate != "" {
			_, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				log.Fatalf("Invalid end date format. Please use YYYY-MM-DD.")
			}
			listTransactionParams.End = &endDate
		}

		var filterAccounts []string
		if filterAccountsStr != "" {
			filterAccounts = strings.Split(filterAccountsStr, ",")
		}

		config := mercury.Config{
			APIKey: apiKey,
		}

		accountsResponse, err := mercury.ListAccounts(config)
		if err != nil {
			log.Fatalf("Something went wrong listing accounts: %v", err)
		}

		accounts := accountsResponse.Accounts
		accounts = filter(accounts, func(account mercury.Account) bool {
			return filterAccountsStr == "" || contains(filterAccounts, account.ID)
		})

		var transactionResponse *mercury.TransactionResponse
		var transactions []mercury.Transaction

		for _, account := range accounts {
			transactionResponse, err = mercury.ListTransactions(config, account.ID, *listTransactionParams)
			if err != nil {
				log.Fatalf("Something went wrong listing transactions for account: %v", err)
			}

			transactions = append(transactions, transactionResponse.Transactions...)
		}

		if counterparty != "" {
			counterpartyFilter := strings.Split(counterparty, ",")
			transactions = filter(transactions, func(txn mercury.Transaction) bool {
				return contains(counterpartyFilter, txn.CounterpartyName)
			})
		}

		transactions = filter(transactions, func(txn mercury.Transaction) bool {
			return txn.Amount < 0
		})

		aggregateSpend := aggregateByCounterparty(transactions, aggregate)

		for key, val := range aggregateSpend {
			fmt.Printf("%s: $%.2f\n", key, val)
		}
	},
}

var revenueCmd = &cobra.Command{
	Use:   "revenue",
	Short: "Track revenue across counterparties and periods",
	Run: func(cmd *cobra.Command, args []string) {
		if apiKey == "" {
			apiKey = viper.GetString("MERCURY_API_KEY")
		}

		if apiKey == "" {
			log.Fatalf("You must provide an API key, either by setting MERCURY_API_KEY or passing --api-key")
		}

		listTransactionParams := &mercury.ListTransactionsParams{}

		if startDate != "" {
			_, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				log.Fatalf("Invalid start date format. Please use YYYY-MM-DD.")
			}
			listTransactionParams.Start = &startDate
		}

		if endDate != "" {
			_, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				log.Fatalf("Invalid end date format. Please use YYYY-MM-DD.")
			}
			listTransactionParams.End = &endDate
		}

		var filterAccounts []string
		if filterAccountsStr != "" {
			filterAccounts = strings.Split(filterAccountsStr, ",")
		}

		config := mercury.Config{
			APIKey: apiKey,
		}

		accountsResponse, err := mercury.ListAccounts(config)
		if err != nil {
			log.Fatalf("Something went wrong listing accounts: %v", err)
		}

		accounts := accountsResponse.Accounts
		accounts = filter(accounts, func(account mercury.Account) bool {
			return filterAccountsStr == "" || contains(filterAccounts, account.ID)
		})

		var transactionResponse *mercury.TransactionResponse
		var transactions []mercury.Transaction

		for _, account := range accounts {
			transactionResponse, err = mercury.ListTransactions(config, account.ID, *listTransactionParams)
			if err != nil {
				log.Fatalf("Something went wrong listing transactions for account: %v", err)
			}

			transactions = append(transactions, transactionResponse.Transactions...)
		}

		if counterparty != "" {
			counterpartyFilter := strings.Split(counterparty, ",")
			transactions = filter(transactions, func(txn mercury.Transaction) bool {
				return contains(counterpartyFilter, txn.CounterpartyName)
			})
		}

		transactions = filter(transactions, func(txn mercury.Transaction) bool {
			return txn.Amount > 0
		})

		aggregateSpend := aggregateByCounterparty(transactions, aggregate)

		for key, val := range aggregateSpend {
			fmt.Printf("%s: $%.2f\n", key, val)
		}
	},
}

func init() {
	rootCmd.AddCommand(spendCmd)
	rootCmd.AddCommand(balancesCmd)
	rootCmd.AddCommand(revenueCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	balancesCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	balancesCmd.Flags().StringVar(&filterAccountsStr, "accounts", "", "Specifies accounts to filter")

	spendCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	spendCmd.Flags().StringVar(&filterAccountsStr, "accounts", "", "Specifies accounts to filter")
	spendCmd.Flags().StringVar(&startDate, "start-date", "", "Specifies the start date for tracking expenses (YYYY-MM-DD)")
	spendCmd.Flags().StringVar(&endDate, "end-date", "", "Specifies the end date for tracking expenses (YYYY-MM-DD)")
	spendCmd.Flags().StringVar(&aggregate, "aggregate", "counterparty", "Specifies how the expenses should be broken down (counterparty, day, month, year, multiple)")
	spendCmd.Flags().StringVar(&counterparty, "counterparty", "", "Filters expenses for specific counterparties")

	revenueCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	revenueCmd.Flags().StringVar(&filterAccountsStr, "accounts", "", "Specifies accounts to filter")
	revenueCmd.Flags().StringVar(&startDate, "start-date", "", "Specifies the start date for tracking revenue (YYYY-MM-DD)")
	revenueCmd.Flags().StringVar(&endDate, "end-date", "", "Specifies the end date for tracking revenue (YYYY-MM-DD)")
	revenueCmd.Flags().StringVar(&aggregate, "aggregate", "counterparty", "Specifies how the revenue should be broken down (counterparty, day, month, year, multiple)")
	revenueCmd.Flags().StringVar(&counterparty, "counterparty", "", "Filters revenue for specific counterparties")

	// Bind environment variables
	viper.AutomaticEnv()
	viper.BindEnv("MERCURY_API_KEY")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func contains(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

func aggregateByCounterparty(txns []mercury.Transaction, aggregate string) map[string]float64 {
	// TODO: Use aggregation
	spend := make(map[string]float64)

	for _, txn := range txns {
		if currentSpend, ok := spend[txn.CounterpartyName]; ok {
			spend[txn.CounterpartyName] = txn.Amount + currentSpend
		} else {
			spend[txn.CounterpartyName] = txn.Amount
		}
	}

	return spend
}
