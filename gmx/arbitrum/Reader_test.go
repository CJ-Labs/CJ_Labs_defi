package arbitrum

import (
	"context"
	"github.com/CJ-Labs/CJ_Labs_defi/event/arbitrum/reader"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestGetMarketInfo(t *testing.T) {
	// 1. 初始化以太坊客户端
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID") // 替换为你的 Infura 或其他节点 URL
	assert.NoError(t, err, "Failed to connect to Ethereum client")

	// 2. 创建 ReaderEventManager 合约实例
	contractAddress := common.HexToAddress("0x0537C767cDAC0726c76Bb89e92904fe28fd02fE1") // 替换为实际的合约地址
	readerContract, err := reader.NewReaderEventManager(contractAddress, client)
	assert.NoError(t, err, "Failed to instantiate ReaderEventManager contract")

	// 3. 准备测试参数
	dataStoreAddress := common.HexToAddress("0xFD70de6b91282D8017aA4E741e9Ae325CAb992d8") // 替换为实际的 DataStore 合约地址
	marketKey := common.HexToAddress("0xMARKET_ADDRESS")                                  // 替换为你要查询的市场地址

	// 创建市场价格结构体
	prices := reader.MarketUtilsMarketPrices{
		IndexTokenPrice: reader.PriceProps{
			Min: big.NewInt(1000000000000000000), // 1 ETH
			Max: big.NewInt(1010000000000000000), // 1.01 ETH
		},
		LongTokenPrice: reader.PriceProps{
			Min: big.NewInt(1000000000000000000), // 1 ETH
			Max: big.NewInt(1010000000000000000), // 1.01 ETH
		},
		ShortTokenPrice: reader.PriceProps{
			Min: big.NewInt(1000000000000000000), // 1 USDC (假设 1:1)
			Max: big.NewInt(1000000000000000000), // 1 USDC
		},
	}

	// 4. 调用合约方法
	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}

	marketInfo, err := readerContract.GetMarketInfo(callOpts, dataStoreAddress, prices, marketKey)
	assert.NoError(t, err, "Failed to get market info")

	// 5. 验证返回结果
	t.Logf("Market Info: %+v", marketInfo)

	// 检查市场基本信息
	assert.NotEqual(t, common.Address{}, marketInfo.Market.MarketToken, "Market token address should not be empty")
	assert.NotEqual(t, common.Address{}, marketInfo.Market.IndexToken, "Index token address should not be empty")
	assert.NotEqual(t, common.Address{}, marketInfo.Market.LongToken, "Long token address should not be empty")
	assert.NotEqual(t, common.Address{}, marketInfo.Market.ShortToken, "Short token address should not be empty")

	// 检查借贷因子
	assert.NotNil(t, marketInfo.BorrowingFactorPerSecondForLongs, "Borrowing factor for longs should not be nil")
	assert.NotNil(t, marketInfo.BorrowingFactorPerSecondForShorts, "Borrowing factor for shorts should not be nil")

	// 检查基础资金费率
	assert.NotNil(t, marketInfo.BaseFunding.FundingFeeAmountPerSize.Long.LongToken, "Long token funding fee should not be nil")
	assert.NotNil(t, marketInfo.BaseFunding.FundingFeeAmountPerSize.Long.ShortToken, "Short token funding fee for longs should not be nil")
	assert.NotNil(t, marketInfo.BaseFunding.ClaimableFundingAmountPerSize.Short.LongToken, "Long token claimable funding for shorts should not be nil")
	assert.NotNil(t, marketInfo.BaseFunding.ClaimableFundingAmountPerSize.Short.ShortToken, "Short token claimable funding for shorts should not be nil")

	// 检查虚拟库存
	assert.NotNil(t, marketInfo.VirtualInventory.VirtualPoolAmountForLongToken, "Virtual pool amount for long token should not be nil")
	assert.NotNil(t, marketInfo.VirtualInventory.VirtualPoolAmountForShortToken, "Virtual pool amount for short token should not be nil")
	assert.NotNil(t, marketInfo.VirtualInventory.VirtualInventoryForPositions, "Virtual inventory for positions should not be nil")
}

func TestGetMarketInfoList(t *testing.T) {
	// 1. 初始化以太坊客户端
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID") // 替换为你的 Infura 或其他节点 URL
	assert.NoError(t, err, "Failed to connect to Ethereum client")

	// 2. 创建 ReaderEventManager 合约实例
	contractAddress := common.HexToAddress("0xCONTRACT_ADDRESS") // 替换为实际的合约地址
	readerContract, err := reader.NewReaderEventManager(contractAddress, client)
	assert.NoError(t, err, "Failed to instantiate ReaderEventManager contract")

	// 3. 准备测试参数
	dataStoreAddress := common.HexToAddress("0xDATA_STORE_ADDRESS") // 替换为实际的 DataStore 合约地址

	// 创建市场价格列表
	var marketPricesList []reader.MarketUtilsMarketPrices
	for i := 0; i < 3; i++ {
		marketPrices := reader.MarketUtilsMarketPrices{
			IndexTokenPrice: reader.PriceProps{
				Min: big.NewInt(1000000000000000000), // 1 ETH
				Max: big.NewInt(1010000000000000000), // 1.01 ETH
			},
			LongTokenPrice: reader.PriceProps{
				Min: big.NewInt(1000000000000000000), // 1 ETH
				Max: big.NewInt(1010000000000000000), // 1.01 ETH
			},
			ShortTokenPrice: reader.PriceProps{
				Min: big.NewInt(1000000000000000000), // 1 USDC (假设 1:1)
				Max: big.NewInt(1000000000000000000), // 1 USDC
			},
		}
		marketPricesList = append(marketPricesList, marketPrices)
	}

	start := big.NewInt(0)
	end := big.NewInt(3)

	// 4. 调用合约方法
	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}

	marketInfos, err := readerContract.GetMarketInfoList(callOpts, dataStoreAddress, marketPricesList, start, end)
	assert.NoError(t, err, "Failed to get market info list")

	// 5. 验证返回结果
	assert.Equal(t, 3, len(marketInfos), "Should return 3 market infos")

	for i, marketInfo := range marketInfos {
		t.Logf("Market %d Info: %+v", i, marketInfo)

		// 检查每个市场的基本信息
		assert.NotEqual(t, common.Address{}, marketInfo.Market.MarketToken, "Market token address should not be empty")
		assert.NotEqual(t, common.Address{}, marketInfo.Market.IndexToken, "Index token address should not be empty")
		assert.NotEqual(t, common.Address{}, marketInfo.Market.LongToken, "Long token address should not be empty")
		assert.NotEqual(t, common.Address{}, marketInfo.Market.ShortToken, "Short token address should not be empty")
	}
}

func TestGetMarket(t *testing.T) {
	// 1. 初始化以太坊客户端
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID") // 替换为你的 Infura 或其他节点 URL
	assert.NoError(t, err, "Failed to connect to Ethereum client")

	// 2. 创建 ReaderEventManager 合约实例
	contractAddress := common.HexToAddress("0xCONTRACT_ADDRESS") // 替换为实际的合约地址
	readerContract, err := reader.NewReaderEventManager(contractAddress, client)
	assert.NoError(t, err, "Failed to instantiate ReaderEventManager contract")

	// 3. 准备测试参数
	dataStoreAddress := common.HexToAddress("0xDATA_STORE_ADDRESS") // 替换为实际的 DataStore 合约地址
	marketKey := common.HexToAddress("0xMARKET_ADDRESS")            // 替换为你要查询的市场地址

	// 4. 调用合约方法
	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}

	marketProps, err := readerContract.GetMarket(callOpts, dataStoreAddress, marketKey)
	assert.NoError(t, err, "Failed to get market")

	// 5. 验证返回结果
	t.Logf("Market Props: %+v", marketProps)

	assert.NotEqual(t, common.Address{}, marketProps.MarketToken, "Market token address should not be empty")
	assert.NotEqual(t, common.Address{}, marketProps.IndexToken, "Index token address should not be empty")
	assert.NotEqual(t, common.Address{}, marketProps.LongToken, "Long token address should not be empty")
	assert.NotEqual(t, common.Address{}, marketProps.ShortToken, "Short token address should not be empty")
}

func TestGetMarkets(t *testing.T) {
	// 1. 初始化以太坊客户端
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID") // 替换为你的 Infura 或其他节点 URL
	assert.NoError(t, err, "Failed to connect to Ethereum client")

	// 2. 创建 ReaderEventManager 合约实例
	contractAddress := common.HexToAddress("0xCONTRACT_ADDRESS") // 替换为实际的合约地址
	readerContract, err := reader.NewReaderEventManager(contractAddress, client)
	assert.NoError(t, err, "Failed to instantiate ReaderEventManager contract")

	// 3. 准备测试参数
	dataStoreAddress := common.HexToAddress("0xDATA_STORE_ADDRESS") // 替换为实际的 DataStore 合约地址
	start := big.NewInt(0)
	end := big.NewInt(10) // 获取前10个市场

	// 4. 调用合约方法
	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}

	markets, err := readerContract.GetMarkets(callOpts, dataStoreAddress, start, end)
	assert.NoError(t, err, "Failed to get markets")

	// 5. 验证返回结果
	assert.True(t, len(markets) > 0, "Should return at least one market")

	for i, market := range markets {
		t.Logf("Market %d: %+v", i, market)

		// 检查每个市场的基本信息
		assert.NotEqual(t, common.Address{}, market.MarketToken, "Market token address should not be empty")
		assert.NotEqual(t, common.Address{}, market.IndexToken, "Index token address should not be empty")
		assert.NotEqual(t, common.Address{}, market.LongToken, "Long token address should not be empty")
		assert.NotEqual(t, common.Address{}, market.ShortToken, "Short token address should not be empty")
	}
}
