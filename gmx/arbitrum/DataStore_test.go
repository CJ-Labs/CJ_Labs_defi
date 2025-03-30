package arbitrum

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/CJ-Labs/CJ_Labs_defi/event/arbitrum/datastore"
)

// 辅助函数：检查地址是否在列表中
func containsAddress(addrs []common.Address, addr common.Address) bool {
	for _, a := range addrs {
		if a == addr {
			return true
		}
	}
	return false
}

func TestGetWbtcUsdcMarketDataDirectly(t *testing.T) {
	// 1. 连接到 Arbitrum One
	client, err := ethclient.Dial("https://arb1.arbitrum.io/rpc")
	require.NoError(t, err, "Failed to connect to Arbitrum One")
	defer client.Close()

	// 2. GMX v2 DataStore 合约地址
	dataStoreAddress := common.HexToAddress("0xFD70de6b91282D8017aA4E741e9Ae325CAb992d8")
	// marketstoreutils
	//marketstoreutilsAddress := common.HexToAddress("0x1C7AB4104a8E43a5CC0688143EFb284E4045d32c")

	// 3. WBTC/USDC 代币地址
	wbtcAddress := common.HexToAddress("0x2f2a2543B76A4166549F7aaB2e75Bef0aefC5B0f")
	usdcAddress := common.HexToAddress("0xaf88d065e77c8cC2239327C5EDb3A432268e5831")

	// 2. DataStore 合约地址
	dataStore, err := datastore.NewDataStoreEventManager(dataStoreAddress, client)
	if err != nil {
		t.Fatalf("Failed to instantiate DataStore: %v", err)
	}
	require.NoError(t, err)

	//marketStoreUtils, err := marketstoreutils.NewMarketStoreUtilsEventManager(marketstoreutilsAddress, client)
	//if err != nil {
	//	t.Fatalf("Failed to instantiate DataStore: %v", err)
	//}
	//require.NoError(t, err)

	// 4. 定义查询函数
	getMarketUint := func(marketKey, dataKey []byte) (*big.Int, error) {
		fullKey := crypto.Keccak256(marketKey, dataKey)
		return dataStore.GetUint(nil, common.BytesToHash(fullKey))
	}

	t.Run("Get Market address", func(t *testing.T) {
		// 示例值，请替换为实际数据
		indexToken := wbtcAddress      // 指数代币地址
		longToken := wbtcAddress       // 多头代币地址
		shortToken := usdcAddress      // 空头代币地址
		marketType := "someMarketType" // 市场类型

		// 在 Go 中模拟 abi.encode
		stringType, _ := abi.NewType("string", "", nil)   // 定义字符串类型
		addressType, _ := abi.NewType("address", "", nil) // 定义地址类型

		arguments := abi.Arguments{
			{Type: stringType},  // "GMX_MARKET"
			{Type: addressType}, // indexToken
			{Type: addressType}, // longToken
			{Type: addressType}, // shortToken
			{Type: stringType},  // marketType
		}

		// 将参数打包为字节数组，类似 Solidity 的 abi.encode
		encoded, err := arguments.Pack("GMX_MARKET", indexToken, longToken, shortToken, marketType)
		if err != nil {
			panic(err) // 如果打包失败，抛出异常
		}

		// 计算 keccak256 哈希
		hash := sha3.NewLegacyKeccak256() // 创建 keccak256 哈希对象
		hash.Write(encoded)               // 写入编码后的数据
		salt := hash.Sum(nil)             // 获取哈希结果

		var saltBytes32 [32]byte        // 定义 bytes32 类型
		copy(saltBytes32[:], salt[:32]) // 将哈希结果复制到定长数组

		//// 检查市场是否已存在
		//if existingMarketAddress != common.HexToAddress("0x0000000000000000000000000000000000000000") {
		//	panic(fmt.Sprintf("市场已存在: salt=%s, existingMarketAddress=%s",
		//		hex.EncodeToString(saltBytes32[:]), existingMarketAddress.Hex()))
		//}

	})

	// 5. 查询多头总头寸
	t.Run("Get Long Positions", func(t *testing.T) {
		// 构造市场键 (isLong = true)
		marketKey := crypto.Keccak256(
			common.LeftPadBytes(wbtcAddress.Bytes(), 32),
			common.LeftPadBytes(usdcAddress.Bytes(), 32),
			[]byte{1}, // isLong = true
		)

		// 查询多头总头寸
		totalLong, err := getMarketUint(marketKey, []byte("TOTAL_POSITION_AMOUNT"))
		require.NoError(t, err)
		t.Logf("WBTC/USDC Total Long Positions: %s", totalLong)
	})

	// 6. 查询空头总头寸
	t.Run("Get Short Positions", func(t *testing.T) {
		// 构造市场键 (isLong = false)
		marketKey := crypto.Keccak256(
			common.LeftPadBytes(wbtcAddress.Bytes(), 32),
			common.LeftPadBytes(usdcAddress.Bytes(), 32),
			[]byte{0}, // isLong = false
		)

		// 查询空头总头寸
		totalShort, err := getMarketUint(marketKey, []byte("TOTAL_POSITION_AMOUNT"))
		require.NoError(t, err)
		t.Logf("WBTC/USDC Total Short Positions: %s", totalShort)
	})

	// 7. 查询其他市场数据示例
	t.Run("Get Market Status", func(t *testing.T) {
		// 构造市场键 (不包含isLong)
		marketKey := crypto.Keccak256(
			common.LeftPadBytes(wbtcAddress.Bytes(), 32),
			common.LeftPadBytes(usdcAddress.Bytes(), 32),
		)

		// 查询市场是否启用
		isDisabled, err := getMarketUint(marketKey, []byte("IS_MARKET_DISABLED"))
		require.NoError(t, err)
		t.Logf("Is Market Disabled: %s (0 = enabled, 1 = disabled)", isDisabled)
	})

}
