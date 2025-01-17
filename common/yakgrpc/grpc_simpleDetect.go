package yakgrpc

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const simpleDetect = `
yakit.AutoInitYakit()

loglevel("info")

runtimeID = env.Get("YAK_RUNTIME_ID")
yakit.StatusCard("RuntimeIDFromRisks", runtimeID)

# Input your code!
updateStatus = s => yakit.StatusCard("当前状态", s)

updateStatus("解析参数")
yakit.Info("解析参数")
yakit.Info("开始处理扫描参数")

targetFiles = cli.String("target-file")

hostCount = cli.Int("host-count")

// 恢复任务的指针
recordPtr = cli.Int("record-ptr", cli.setDefault(0))

progressData = cli.Float("progress-data")

baseProgress = 0.05


if progressData > 0 {
	baseProgress = progressData
}
yakit.SetProgress(baseProgress)

taskName = cli.String("task-name")

tcpPorts = cli.String("ports", cli.setDefault("22,443,445,80,8000-8004,3306,3389,5432,8080-8084,7000-7005"))
mode = cli.String("mode", cli.setDefault("fingerprint"))
saveToDB = cli.Bool("save-to-db")
saveClosed = cli.Bool("save-closed-ports") 
proxies = cli.String("proxy", cli.setDefault("no"))
probeTimeoutFloat = cli.Float("probe-timeout", cli.setDefault(5.0), cli.setRequired(false))
probeMax = cli.Int("probe-max", cli.setRequired(false), cli.setDefault(4))

// host alive scan
skippedHostAliveScan = cli.Bool("skipped-host-alive-scan")
hostAliveConcurrent = cli.Int("host-alive-concurrent", cli.setDefault(20), cli.setRequired(false))
hostAliveTimeout = cli.Float("host-alive-timeout", cli.setDefault(5.0), cli.setRequired(false))
hostAliveTCPPorts = cli.String("host-alive-ports", cli.setDefault("80,22,443"), cli.setRequired(false))

if proxies == "no" {
	proxies = ""
}


active = cli.Bool("active")
concurrent = cli.Int("concurrent", cli.setDefault(50))
synConcurrent = cli.Int("syn-concurrent", cli.setDefault(800))
protos = cli.String("proto", cli.setDefault("tcp"))

fpMode = cli.String("fp-mode", cli.setDefault("all"))
scriptNameFile = cli.String("script-name-file", cli.setDefault(""))

cli.check()
yakit.SetProgress(baseProgress+0.01)

excludeHosts = cli.String("exclude-hosts")
excludePorts = cli.String("exclude-ports")
enableCrawler = cli.Have("enable-basic-crawler")
crawlerMaxRequest = cli.Int("basic-crawler-request-max")

tcpOps = []
tcpOps = append(tcpOps, servicescan.active(active))

if concurrent > 0 {
    tcpOps = append(tcpOps, servicescan.concurrent(concurrent))
}

if protos != "" {
    protoList = str.Split(protos, ",")
	printf("PROTO: %#v\n", protos)
    tcpOps = append(tcpOps, servicescan.proto(protoList...))
}

// 使用指纹检测规则条数
if probeMax > 0 {
	tcpOps = append(tcpOps, servicescan.maxProbes(probeMax))
} else {
	tcpOps = append(tcpOps, servicescan.maxProbes(3))
}

if proxies != "" {
    proxyList = str.Split(proxies, ",")
    printf("PROXY: %v\n", proxyList)
    tcpOps = append(tcpOps, servicescan.proxy(proxyList...))
}

if probeTimeoutFloat > 0 {
    tcpOps = append(tcpOps, servicescan.probeTimeout(probeTimeoutFloat))
}

if fpMode == "web" {
	tcpOps = append(tcpOps, servicescan.web())
}

if fpMode == "service" {
	tcpOps = append(tcpOps, servicescan.service())
}

if fpMode == "all" {
	tcpOps = append(tcpOps, servicescan.all())
}

// 设置排除端口和主机
if excludePorts != "" {
    yakit.Info("设置排除端口：%v", excludePorts)
    tcpOps = append(tcpOps, servicescan.excludePorts(excludePorts))
}
if excludeHosts != "" {
    yakit.Info("设置排除主机：%v", excludeHosts)
    tcpOps = append(tcpOps, servicescan.excludeHosts(excludeHosts))
}

// top 5000 ports
synPorts = "1-4,6-17,19-35,37-38,40,42-45,47-60,63,65-104,106,108-120,122-125,127-130,132-133,135-139,141-146,148-151,156-158,161-163,165,168,173-174,176-177,179-182,184-185,189-194,196,199-202,204-206,209-214,216-217,219-226,228-231,233-238,243,247-262,264-265,267-268,270-271,273,276-277,280,284,288-289,293-295,300-301,303,305-306,308,311,315-316,322,325-326,329,333-334,336-337,340,343,346,350-353,355,358,360-362,364,366,369-370,372-373,379-380,383,388-389,391-392,395,397,399-404,406-408,410-420,422-423,425,427-428,432-435,437-454,456-458,460,462,464-466,470,472-475,479-481,485-487,491-493,496-497,500-502,505,507,509-516,518-519,522-528,530,533-536,538,540-545,548,552-557,560-561,563-564,568-572,577-578,580,582-583,587,590-591,593,596,598-649,651-670,672-678,680-692,694-771,773-916,918-939,941-1114,1116-1119,1121-1128,1130-1132,1134-1139,1141-1145,1147-1154,1156-1159,1162-1169,1173-1176,1178-1180,1182-1192,1194-1196,1198-1201,1204,1207-1218,1220-1223,1228-1229,1233-1234,1236,1239-1241,1243-1244,1247-1251,1259,1261-1262,1264,1268,1270-1272,1276-1277,1279,1282,1287,1290-1291,1296-1297,1299-1303,1305-1312,1314-1319,1321-1322,1324,1327-1328,1330-1331,1334,1336-1340,1346-1374,1376,1379-1391,1393-1405,1407-1414,1416-1420,1422-1427,1429-1430,1432-1446,1448-1451,1453-1459,1461-1462,1464-1467,1469-1470,1472-1476,1479-1480,1482-1484,1486,1488,1491-1503,1505-1511,1513,1515-1519,1521-1529,1531-1533,1535,1537-1545,1547-1552,1556,1558-1560,1565-1566,1569,1580,1583-1584,1592,1594,1598,1600,1605,1607,1615,1620,1622,1627,1631-1632,1635,1638,1641,1645,1650-1652,1658,1661-1664,1666-1668,1670-1672,1677,1680,1683,1687-1688,1691,1694,1699-1701,1703,1707-1709,1711-1713,1715,1717-1723,1730,1735-1736,1745,1750,1752-1753,1755,1759,1761-1763,1782-1783,1785,1791-1792,1796,1799-1801,1805-1808,1811-1812,1814,1823,1825,1827,1832,1835,1839-1840,1846,1852,1858,1860-1864,1870-1871,1875,1883,1885,1889,1900-1901,1906,1910-1912,1914,1918,1924-1925,1927,1935,1947,1954,1956,1958,1964-1965,1967,1971-1977,1979,1981,1984,1986-2016,2018-2028,2030-2031,2033-2035,2038,2040-2049,2053-2054,2060,2062,2064-2065,2067-2070,2074,2080-2084,2086-2087,2095-2096,2099-2108,2111-2112,2115,2119-2121,2124,2126,2132,2134-2135,2142,2144,2148,2150,2154,2156-2158,2160-2161,2165,2170,2179,2187,2190-2191,2196-2197,2200-2201,2203,2216,2222,2224,2231-2232,2241,2248,2250-2251,2253,2255,2260-2262,2265,2269-2271,2280,2288,2290-2292,2296,2300-2302,2304,2306-2307,2312-2313,2323,2325-2326,2330,2335,2340-2341,2357,2366,2371-2372,2375-2376,2381-2383,2386-2387,2391-2394,2399,2401-2403,2405,2409-2410,2418,2425-2426,2430-2433,2435-2436,2438-2440,2449,2456,2463,2470,2472,2474,2492-2493,2500-2501,2505,2522,2525,2529-2532,2537,2539,2542,2550-2551,2554,2557-2558,2564,2567-2569,2573,2576,2580,2583-2584,2595-2596,2598-2602,2604-2610,2622-2623,2627-2628,2631,2638,2644,2655,2657-2659,2670,2675,2679,2691,2696,2700-2704,2706,2710-2712,2717-2718,2721,2723,2725,2728-2729,2734,2737,2739,2766,2771,2784,2800,2802,2804-2806,2809,2811-2812,2821,2847,2850,2862-2863,2869,2875-2876,2880,2882,2888-2889,2898,2901-2903,2907-2910,2915,2920,2930,2948-2949,2957-2958,2967-2968,2973,2977-2978,2984,2987-2989,2991,2993-2994,2997-2998,3000-3003,3005-3007,3011-3014,3017-3018,3023,3025,3030-3031,3040,3045,3049-3050,3052,3057,3061-3064,3071,3073-3074,3077,3080,3086,3089,3098,3102-3103,3111,3118-3119,3121,3128,3141,3146,3158,3160,3162,3167-3168,3175,3185,3190,3200,3207,3210-3212,3217,3219-3221,3230-3231,3233,3240-3241,3243,3245,3249,3251,3260-3261,3263-3264,3268-3269,3280-3281,3283,3291-3292,3299-3301,3304,3306-3311,3313-3314,3317-3319,3322-3325,3333-3334,3346,3351,3359,3362-3363,3365,3367-3372,3374-3376,3380,3388-3390,3396-3400,3404,3410,3413-3415,3419,3421,3423-3425,3427,3430,3438-3439,3443,3456-3457,3459,3468,3474,3476,3478-3479,3482-3483,3485-3486,3491-3493,3495,3497,3503,3505-3506,3509-3511,3513-3515,3517,3519-3521,3526-3527,3530-3532,3534,3541,3544,3546-3547,3549,3551-3552,3555,3563,3577,3579-3580,3584,3586,3590,3599-3603,3606,3608,3615,3621-3622,3632-3633,3636-3637,3646-3647,3651-3654,3656,3658-3659,3663-3664,3669-3670,3672,3674,3676,3680-3681,3683-3684,3686,3689-3690,3696-3697,3700,3702-3703,3708,3712-3713,3716-3717,3720,3722-3724,3726,3728,3731,3734,3737-3738,3741-3742,3749,3752,3754,3763,3765-3766,3784,3786-3788,3790-3793,3795-3796,3798-3801,3803,3805-3814,3816-3817,3820,3823-3828,3830-3831,3834-3835,3837,3839-3840,3842,3845-3854,3856-3861,3863,3865-3872,3876-3880,3882-3885,3888-3890,3894-3895,3897,3899-3902,3904-3909,3911,3913-3916,3918-3920,3922-3923,3925,3928-3931,3935-3937,3940-3941,3943-3946,3948-3949,3952,3956-3957,3959,3961-3964,3967-3969,3971-3972,3975,3979-3986,3989-4010,4015-4020,4022,4024-4025,4029,4031-4032,4034-4036,4039-4040,4045,4049-4050,4052,4056,4058,4065,4070,4072,4080,4087,4090,4094-4096,4100-4101,4111-4113,4118-4121,4123,4125-4126,4129,4132-4133,4135,4141,4143-4147,4156-4158,4161,4164-4165,4170,4174,4176-4177,4184-4185,4188,4190,4192,4199-4300,4302,4307,4321-4323,4325,4328,4333,4342-4343,4355-4358,4368-4369,4374-4376,4384,4388-4389,4401,4407,4414-4415,4417-4418,4422-4423,4426,4430,4433,4442-4447,4449,4454,4464,4471,4476,4480,4500,4516-4517,4530,4534,4537-4538,4545,4550,4555,4557-4559,4567,4570,4599-4602,4606,4609,4644,4649,4658,4660,4662,4665,4672,4683,4687,4689,4700,4711-4713,4725,4738,4745,4760,4767,4770-4771,4778,4786,4789,4793,4800,4819,4827,4845,4847-4848,4850,4859-4860,4867,4875-4877,4879,4881-4882,4888-4889,4899-4900,4903,4912,4931,4949-4950,4953,4987,4991,4998-5017,5020-5021,5023,5026-5027,5030,5033,5040,5046,5048,5050-5055,5060-5061,5063,5066,5070,5074,5080-5081,5087-5088,5090,5095-5096,5098,5100-5102,5111,5114-5115,5120-5122,5125,5133,5137,5145-5147,5151-5153,5161-5162,5164,5166,5190-5191,5193,5195-5197,5200-5203,5212,5214,5219,5221-5223,5225-5226,5232-5235,5242,5250,5252,5259,5261,5269,5279-5281,5291,5298,5300-5303,5306,5308,5320-5321,5339,5344,5347,5349,5353,5357-5358,5363,5370,5377,5397,5400,5405-5406,5414,5416-5417,5423,5431-5433,5440-5442,5444,5450,5457-5458,5462-5463,5470-5475,5490,5500-5502,5510,5520,5530,5544,5550,5552-5555,5557,5560,5566,5580-5583,5611-5612,5620-5622,5629,5631-5633,5637-5639,5646,5665-5667,5670,5672,5675,5678-5680,5705,5711,5713-5714,5717-5718,5721-5723,5730,5732,5734,5737,5743-5744,5747-5748,5769-5770,5780,5793,5800-5804,5806-5808,5810-5812,5814-5815,5817-5818,5820-5827,5831,5834,5836,5838-5841,5845,5848-5850,5852-5854,5858-5860,5862,5868-5869,5871,5874-5875,5877-5878,5881,5887-5888,5899-5912,5914-5915,5917-5918,5920-5927,5931,5934,5936,5938-5940,5945,5948-5950,5952-5954,5958-5963,5966,5968-5969,5971,5974-5975,5977-5978,5981,5985-5990,5997-6063,6065,6067-6068,6077,6083,6085,6090-6091,6100-6101,6103,6105-6106,6110-6116,6118,6120-6121,6123,6126,6129-6130,6141-6147,6156,6161,6203,6222,6247,6250-6252,6259,6273-6274,6306,6309-6310,6323-6324,6346-6347,6349-6350,6379,6389,6400-6401,6412,6418-6419,6442,6481,6483,6500,6502-6504,6510,6514,6520,6535,6543-6544,6547-6548,6550-6551,6565-6567,6579-6580,6588,6600,6602,6606,6619,6626,6628,6644,6646-6647,6650,6662,6665-6670,6672-6673,6678,6689,6692,6699-6701,6709-6711,6725,6732,6734,6779-6780,6786-6789,6792,6839,6877,6881,6888,6896-6897,6901,6920,6922,6936,6942,6956,6969,6972-6973,7000-7012,7015-7016,7019,7024-7025,7033,7043,7050-7051,7067-7068,7070-7072,7080,7088,7092,7099-7104,7106,7119,7121,7123,7173,7184,7200-7201,7218,7231,7236,7241,7272-7273,7278,7281,7300-7359,7400,7402,7411,7435,7438,7443,7451,7456,7464,7471,7496,7500-7501,7512,7548,7553,7555,7560,7566,7588,7597,7600,7625,7627-7628,7631,7634,7637,7654,7676-7677,7685,7688,7725,7727,7734,7741,7744,7749,7770-7772,7777-7778,7780,7788-7789,7800,7802,7813,7830,7852-7854,7872,7878,7887,7895,7900-7903,7911,7913,7920-7921,7929,7932-7933,7937-7938,7975,7981-7982,7998-8003,8005-8011,8014-8016,8018-8019,8021-8023,8025,8029,8031,8034,8037,8041-8042,8045,8050,8052,8060,8064,8066,8076,8080-8090,8092-8093,8095,8097-8100,8110,8116,8118,8123,8133,8140,8144,8180-8182,8189,8192-8194,8199-8202,8222,8232,8243,8245,8248,8254-8255,8268,8273,8280,8282,8290-8295,8300,8308,8320-8321,8333,8339,8383,8385,8400-8405,8409,8443,8445,8451-8455,8470-8474,8477,8479,8481,8484,8500,8515,8530-8531,8539-8540,8555,8562,8600-8601,8621,8640,8644,8648-8649,8651-8652,8654-8655,8658,8673,8675-8676,8680,8686,8699,8701,8736,8752,8756,8765-8766,8770,8772,8778,8790,8798,8800-8801,8808,8843,8865,8873,8877-8880,8882-8883,8887-8889,8892,8898-8900,8912,8925,8937,8953-8954,8980,8987,8989-8991,8994,8996,8999-9005,9009-9011,9013,9020-9022,9037,9040,9044,9050-9051,9061,9065,9071,9080-9081,9084,9086,9090-9092,9098-9107,9110-9111,9125,9128,9130-9131,9133,9152,9160-9161,9170,9183,9191,9197-9198,9200-9207,9210-9211,9220,9281-9282,9284,9287,9290,9293,9300,9306,9312,9333,9343,9351,9364,9400,9409,9415,9418,9443-9444,9454,9464,9478,9485,9493,9500-9503,9513,9522,9527,9535-9536,9555,9575,9583,9592-9595,9598,9600,9612-9613,9616,9618-9621,9628-9629,9643,9648,9654,9661,9665-9668,9673-9674,9679-9680,9683,9694,9700,9745,9762,9777,9800,9802,9812,9814-9815,9823,9825-9826,9830,9844,9875-9878,9898,9900-9901,9908-9912,9914-9915,9917,9919,9925,9929,9941,9943-9944,9950,9968,9971,9975,9978-9979,9988,9990-9992,9995,9998-10012,10018-10020,10022-10025,10034-10035,10042,10045,10050-10051,10058,10064,10082-10083,10093,10100-10101,10104,10107,10115,10160-10162,10180,10200,10215,10238,10243,10245-10246,10255,10261,10280,10338,10347,10357,10387,10414,10443,10494,10500,10509,10529,10535,10550-10556,10565-10567,10601-10602,10616-10617,10621,10626,10628-10629,10699,10754,10778,10842,10852,10873,10878,10900,11000-11001,11003,11007,11019,11026,11031-11033,11089,11100,11110-11111,11161-11165,11171,11173,11180,11200,11208,11224,11250,11288,11296,11371,11401,11552,11600,11697,11735,11813,11862-11863,11876-11877,11940,11967,11971,12000-12002,12005-12006,12009,12012-12013,12019,12021,12031,12034,12059,12077,12080,12090,12096-12097,12121,12132,12137,12146,12156,12171,12174,12192,12215,12225,12240,12243,12251,12262,12265,12271,12275,12296,12321-12322,12340,12345-12346,12380,12414,12452,12699,12702,12753,12766,12865,12891-12892,12955,12962,13017,13093,13130,13132,13140,13142,13149,13167,13188,13192-13194,13229,13250,13261,13264-13265,13306,13318,13340,13359,13456,13502,13580,13695,13701,13705,13708-13710,13712-13718,13720-13724,13730,13766,13782-13784,13786,13846,13894,13899,14000-14001,14141,14147,14149-14150,14154,14218,14237-14238,14254,14418,14441-14444,14500,14534,14545,14693,14733,14827,14891,14916,15000-15005,15050,15118,15145,15151,15190-15191,15275,15317,15344-15345,15402,15448,15550,15631,15645-15646,15660,15670,15677,15722,15730,15742,15758,15915,16000-16001,16012,16016,16018,16048,16080,16113,16161-16162,16270,16273,16283,16286,16297,16349,16372,16444,16464,16619,16666,16705,16723-16725,16789,16797,16800,16845,16851,16900-16901,16959,16992-16993,17007,17016-17017,17070,17089,17129,17184-17185,17224-17225,17235,17251,17255,17300,17409,17413,17500,17595,17700-17702,17715,17755-17756,17777,17801-17802,17860,17867,17877,17969,17985,17988,17997,18000,18012,18015,18018,18040,18080,18101,18136,18148,18181-18184,18187,18231,18264,18333,18336-18337,18380,18439,18505,18517,18569,18668-18669,18874,18887,18910,18962,18988,19010,19101,19130,19150,19194,19200-19201,19283,19315,19333,19350,19353,19403,19464,19501,19612,19634,19715,19780,19801,19842,19852,19900,19995-19996,20000-20002,20005,20011,20017,20021,20031-20032,20039,20046,20052,20076,20080,20085,20089,20102,20106,20111,20118,20125,20127,20147,20167,20179-20180,20221-20228,20280,20473,20734,20828,20883,20934,20940,20990,21011,21078,21201,21473,21571,21590,21631,21634,21728,21792,21800,21891,21915,22022,22063,22100,22125,22128,22177,22200,22222-22223,22273,22289-22290,22321,22341,22347,22350-22351,22555,22563,22711,22719,22727,22763,22769,22800,22882,22939,22951,22959,22969,23017,23040,23052,23219,23228,23270,23296,23342,23382,23430,23451,23502,23723,23796,23887,23953,24218,24249,24322-24323,24392,24416,24444,24465,24552,24554,24616,24678,24680,24800,24999-25001,25174,25260,25262,25288,25327,25445,25473,25486,25565,25703,25717,25734-25735,25793,25847,25900,26000-26001,26007,26133,26208,26214,26340,26417,26470,26669,26972,27000-27003,27005,27007,27009-27010,27015-27019,27055,27074-27075,27087,27204,27316,27350-27353,27355-27357,27372,27374,27521,27537,27665,27715,27770,27999,28017,28114,28142,28200-28201,28211,28374,28567,28717,28850-28851,28924,28967,29045,29152,29243,29507,29672,29810,29831,30000-30001,30005,30087,30195,30299,30519,30599,30644,30659,30704-30705,30718,30896,30951,31029,31033,31038,31058,31072,31337,31339,31386,31416,31438,31457,31522,31657,31727-31728,32006,32022,32031,32088,32102,32200,32219,32249,32260-32261,32764-32765,32767-32792,32797-32799,32803,32807,32814-32816,32820,32822,32835,32837,32842,32858,32868-32869,32871,32888,32897-32898,32904-32905,32908,32910-32911,32932,32944,32960-32961,32976,33000,33011,33017,33070,33087,33124,33175,33192,33200,33203,33277,33327,33334-33335,33337,33354,33367,33395,33434,33444,33453,33522-33523,33550,33554,33604-33605,33656,33841,33879,33882,33889,33895,33899,34021,34036,34096,34189,34249,34317,34341,34381,34401,34507,34510,34571-34573,34683,34728,34765,34783,34833,34875,35033,35050,35116,35131,35217,35272,35349,35392-35393,35401,35500,35506,35513,35553,35593,35731,35879,35900-35901,35906,35929,35986,36046,36104-36105,36256,36275,36368,36411,36436,36462,36508,36530,36552,36659,36677,36694,36710,36748,36823-36824,36914,36950,36962,36983,37121,37151,37174,37185,37218,37393,37522,37607,37614,37647,37654,37674,37777,37789,37839,37855,38029,38037,38185,38188,38194,38205,38224,38270,38292,38313,38331,38358,38422,38446,38481,38546,38561,38570,38761,38764,38780,38800,38805,38936,39067,39117,39136,39265,39293,39376,39380,39433,39482,39489,39630,39659,39681,39732,39763,39774,39795,39869,39883,39895,39917,40000-40003,40005,40011,40193,40306,40393,40400,40457,40489,40513,40614,40628,40712,40732,40754,40811-40812,40834,40911,40951,41064,41121,41123,41142,41230,41250,41281,41318,41342,41345,41348,41398,41442,41511,41523,41551,41632,41773,41794-41795,41808,42001,42035,42127,42158,42251,42276,42322,42449,42452,42510,42559-42560,42575,42590,42632,42675,42679,42685,42735,42906,42990,43000,43002,43018,43027,43103,43139,43143,43188,43212,43231,43242,43425,43654,43690,43734,43823,43868,44004,44101,44119,44123,44176,44200,44334,44380,44410,44431,44442-44443,44479,44501,44505,44541,44616,44628,44704,44709,44711,44965,44981,45038,45050,45100,45136,45164,45220,45226,45413,45438,45463,45602,45624,45697,45777,45864,45960,45966,46034,46069,46115,46171,46182,46200,46310,46372,46418,46436,46593,46813,46992,46996,47001,47012,47029,47119,47197,47267,47348,47372,47448,47544,47557,47567,47581,47595,47624,47634,47700,47777,47806,47850,47858,47860,47966,47969,48009,48050,48067,48080,48083,48127,48153,48167,48356,48434,48619,48631,48648,48682,48783,48813,48925,48966-48967,48973,49002,49048,49132,49152-49161,49163-49173,49175-49176,49179,49186,49189-49191,49195-49197,49201-49204,49211,49213,49216,49228,49232,49235-49236,49241,49275,49302,49352,49372,49398,49400-49401,49452,49498,49500,49519-49522,49597,49603,49678,49751,49762,49765,49803,49927,49999-50003,50006,50016,50019,50040,50050,50101,50189,50198,50202,50205,50224,50246,50258,50277,50300,50356,50389,50500,50513,50529,50545,50576-50577,50585,50636,50692,50733,50787,50800,50809,50815,50831,50833-50836,50849,50854,50887,50903,50945,50997,51011,51020,51037,51067,51103,51118,51139,51191,51233-51235,51240,51300,51343,51351,51366,51413,51423,51460,51484-51485,51488,51493,51515,51582,51658,51771-51772,51800,51809,51906,51909,51961,51965,52000-52003,52025,52046,52071,52173,52225-52226,52230,52237,52262,52391,52477,52506,52573,52660,52665,52673,52675,52710,52735,52822,52847-52851,52853,52869,52893,52948,53085,53178,53189,53211-53212,53240,53313-53314,53319,53361,53370,53460,53469,53491,53535,53633,53639,53656,53690,53742,53782,53827,53852,53910,53958,54045,54075,54101,54127,54235,54263,54276,54320-54321,54323,54328,54514,54551,54605,54658,54688,54722,54741,54873,54907,54987,54991,55000,55020,55055-55056,55183,55187,55227,55312,55350,55382,55400,55426,55479,55527,55555-55556,55568-55569,55576,55579,55600,55635,55652,55684,55721,55758,55773,55781,55901,55907,55910,55948,56016,56055,56259,56293,56507,56535,56591,56668,56681,56723,56725,56737-56738,56810,56822,56827,56973,56975,57020,57103,57123,57294,57325,57335,57347,57350,57352,57387,57398,57479,57576,57665,57678,57681,57702,57730,57733,57797,57891,57896,57923,57928,57988,57999,58001-58002,58072,58080,58107,58109,58164,58252,58305,58310,58374,58430,58446,58456,58468,58498,58562,58570,58610,58622,58630,58632,58634,58699,58721,58838,58908,58970,58991,59087,59107,59110,59122,59149,59160,59191,59200-59202,59239,59340,59499,59504,59509-59510,59525,59565,59684,59778,59810,59829,59841,59987,60000,60002-60003,60020,60055,60086,60111,60123,60146,60177,60227,60243,60279,60377,60401,60403,60443,60485,60492,60504,60544,60579,60612,60621,60628,60642,60713,60728,60743,60753,60782-60783,60789,60794,60989,61159,61169-61170,61402,61473,61516,61532,61613,61616-61617,61669,61722,61734,61827,61851,61900,61942,62006,62042,62078,62080,62188,62312,62519,62570,62674,62866,63105,63156,63331,63423,63675,63803,64080,64127,64320,64438,64507,64551,64623,64680,64726-64727,64890,65000,65048,65129,65301,65310-65311,65389,65488,65514"
synPortsList = str.ParseStringToPorts(synPorts)
tcpPortsList = str.ParseStringToPorts(tcpPorts)
strFilter = str.NewFilter()
x.Map(tcpPortsList, p=>{
   strFilter.Insert(f"${p}")
})
newSynPortsList = []
x.Map(synPortsList, p=>{
   if !strFilter.Exist(f"${p}"){
       newSynPortsList.Append(f"${p}")
   }
})
synPorts = str.Join(newSynPortsList, ",")

yakit.SetProgress(baseProgress + 0.02)

var synscanEnable = false
try {
    yakit.Info("检测 SYN 扫描是否可用中")
	yakit.SetProgress(0.1)
    for res in synscan.Scan("127.0.0.1", "80", synscan.wait(1))~ {  }
    
    synscanEnable = true
    yakit.StatusCard("SYN 扫描", "可用")
} catch err {
    yakit.Error("SYN 扫描不可用，原因是：%v", err)
	if str.Contains(sprintf("%v",err),"couldn't load wpcap.dll") {
        yakit.StatusCard("SYN失败", "请安装Npcap", "SYN扫描失败")
    }
}

/*
    插件处理模块
*/
yakit.Info("开始创建插件调用模块")
scriptNames = x.If(scriptNameFile != "", x.Filter(
    x.Map(
        str.ParseStringToLines(string(file.ReadFile(scriptNameFile)[0])), 
        func(e){return str.TrimSpace(e)},
    ), func(e){return e!=""}), make([]string))

pluginsL = len(scriptNames)

yakit.StatusCard("加载插件",f"0/${pluginsL}")

manager, err := hook.NewMixPluginCaller()

if err != nil {
    updateStatus("创建插件调用模块失败")
    die(err)
}
// 这个有必要设置：独立上下文，避免在大并发的时候出现问题
manager.SetConcurrent(20)
manager.SetDividedContext(true)

loadPluginFailed = 0
loadPluginFinished = 0

x.Foreach(scriptNames, func(name){
    err = manager.LoadPlugin(name)
    if err != nil {
        loadPluginFailed++
        yakit.Error("Load Plugin Failed: %s", err)
        return
    }
    
    loadPluginFinished++
	yakit.Info("Load Plugin Success: %s", name)
    return
})
if loadPluginFinished > 0 || loadPluginFailed > 0 {
    yakit.StatusCard(
        "加载插件", 
        f"${loadPluginFinished}/${loadPluginFailed+loadPluginFinished}", 
    )
}else {
	yakit.StatusCard("加载插件失败ID", "请检查是否导入了插件", "加载插件失败")
    die("没有插件加载")
}

yakit.Info("扫描参数设置完成，准备扫描~")


pingOpt =[]

if proxies != "" {
    proxyList = str.Split(proxies, ",")
	pingOpt = append(pingOpt, ping.proxy(proxyList...))
}

if skippedHostAliveScan {
	pingOpt = append(pingOpt, ping.skip(skippedHostAliveScan))
}

if hostAliveTimeout > 0 {
	pingOpt = append(pingOpt, ping.timeout(hostAliveTimeout))
}

if hostAliveConcurrent > 0 {
	pingOpt = append(pingOpt, ping.concurrent(hostAliveConcurrent))
}

if len(hostAliveTCPPorts) > 0 {
	pingOpt = append(pingOpt, ping.tcpPingPorts(hostAliveTCPPorts))
}

if excludeHosts != "" {
	pingOpt = append(pingOpt, ping.excludeHosts(excludeHosts))
}

if !skippedHostAliveScan {
	yakit.Info("开始探测存活~")
	hostRaw, _ = file.ReadFile(targetFiles)
	hosts = str.ReplaceAll(string(hostRaw), "\n" , ",")
	
	aliveHostCount = 0
	ScanHostCount = 0
	showAliveHostCount = func() {
		yakit.StatusCard("存活主机数/扫描主机数",f"${aliveHostCount}/${ScanHostCount}")
	}
	aliveHostStr = ""
	os.RemoveAll(targetFiles)~
	for res := range ping.Scan(hosts, pingOpt...) {
		if res.Ok {
			aliveHostCount ++
			yakit.Info("ping onResult: %v", res)
			aliveHostStr += res.IP + "\n"
		}
		ScanHostCount++
		showAliveHostCount()
	}
	// 如果开启存活检测，那么 host 总数应该为存活的IP
	hostCount = aliveHostCount
	yakit.StatusCard("Ping存活主机数", aliveHostCount)
	aliveHostStr = str.Trim(aliveHostStr, "\n")
	
	file.Save(targetFiles,aliveHostStr)~

	if len(aliveHostStr) == 0 {
		yakit.Error("没有存活主机")
		//die("没有存活主机")
	}
}else{
	yakit.Info("跳过存活探测~")
}



var files = make([]string)

inputFiles = str.SplitAndTrim(targetFiles, ",")
if len(inputFiles) > 0 {
    files.Push(inputFiles...)
}

// updateStatus("解析参数成功")


hostTotal = hostCount
portTotal = len(str.ParseStringToPorts(str.Trim(tcpPorts, ",") + str.Trim(synPorts, ",")))
// 需要传递给报告的参数
yakit.StatusCard("扫描主机数", hostTotal)
yakit.StatusCard("扫描端口数", portTotal)
yakit.StatusCard("开放端口数/已扫主机数", "-/-")


totalTasks = hostTotal * portTotal
portResultFinalTotal = 0
progressLock = sync.NewLock()
updateProgress = func(delta) {
    if totalTasks <= 0 {
        return

    }
	progressLock.Lock()
    defer progressLock.Unlock()

    portResultFinalTotal = portResultFinalTotal+delta
    if portResultFinalTotal > totalTasks {
        portResultFinalTotal = totalTasks
    }
    yakit.Info("已出结果 %v", sprintf("%v/%v/%.3f", portResultFinalTotal, totalTasks,0.1 + (float(portResultFinalTotal) / float(totalTasks) ) * 0.9))
    yakit.SetProgress(baseProgress + 0.03 + (float(portResultFinalTotal) / float(totalTasks) ) * 0.9)
}


handleServiceScanResult = func(result) {
	manager.HandleServiceScanResult(result)
}

/*
    端口扫描部分分为两大块儿
    1. 重点端口 TCP 扫描
    2. 同步非重点端口 SYN 扫描
*/

ScanHostCount = 0
aliveHostCountList = []
scanHostLock = sync.NewLock()

addScanHostCount = func() {
	scanHostLock.Lock()
    defer scanHostLock.Unlock()
    ScanHostCount++
}

OpenPortCount = 0
addOpenPortCount = func() {
    OpenPortCount++
    yakit.StatusCard("开放端口数/已扫主机数", f"${OpenPortCount}/${ScanHostCount}")
	if skippedHostAliveScan {
		aliveHostCount = len(aliveHostCountList)
		yakit.StatusCard("存活主机数/扫描主机数",f"${aliveHostCount}/${ScanHostCount}")
	}
}

/* 进行 TCP 扫描*/

func handleTCP(target) {
    try {
        for result in servicescan.Scan(target, tcpPorts, tcpOps...)~ {
            if !result.IsOpen() {
                continue
            }
			if result.Target not in aliveHostCountList {
				aliveHostCountList = append(aliveHostCountList,result.Target)
			}
            yakit.Info("GOT: %v", result.String())
            yakit.Output({"host": "TCP-"+result.Target, "port": result.Port, "fingerprint": result.GetServiceName(), "htmlTitle": "" + result.GetHtmlTitle(), "isOpen": true})
			yakit.SavePortFromResult(result,taskName)
			addOpenPortCount()
			handleServiceScanResult(result)
        }
        
    } catch err {
        yakit.Error("处理 TCP 指纹识别失败: %v", err)
    } finally {
		p = 1 * len(str.ParseStringToPorts(str.Trim(tcpPorts, ",")))
		updateProgress(p)
	}
}

synWg = sync.NewSizedWaitGroup(2)

func handleSYN(target){
    synWg.Add()
    go func(ctarget) {
        defer synWg.Done()
        try {
            for result in servicescan.ScanFromSynResult(
                synscan.Scan(ctarget, synPorts, 
							synscan.initPortFilter(synPorts),
							synscan.initHostFilter(ctarget),
							synscan.excludePorts(tcpPorts),
							synscan.concurrent(synConcurrent),
					)~, 
					tcpOps..., 
				)~ {
				
                if !result.IsOpen() {
                    continue
                }
                if result.Target not in aliveHostCountList {
					aliveHostCountList = append(aliveHostCountList,result.Target)
				}
                yakit.Info("SYNGOT: %v", result.String())
                yakit.Output({"host": "SYN-"+result.Target, "port": result.Port, "fingerprint": result.GetServiceName(), "htmlTitle": "" + result.GetHtmlTitle(), "isOpen": true})
                yakit.SavePortFromResult(result,taskName)
				
				addOpenPortCount()
				handleServiceScanResult(result)
            }
            
        } catch err { 
            yakit.Error("SYN 处理 TCP 指纹识别失败: %v", err)
        } finally {
    		p = 1 * len(str.ParseStringToPorts(str.Trim(synPorts, ",")))
			updateProgress(p)
		}
    }(target)
}

swg = sync.NewSizedWaitGroup(10)

fileReader, err = file.NewMultiFileLineReader(files...)
if err != nil {
    yakit.Error("无法创建文件读取管理器：%v", err)
    die(err)
}

if recordPtr > 0 {
	yakit.Info("恢复文件指针到 %d", recordPtr)
	err = fileReader.SetRecoverPtr(targetFiles,recordPtr)
	if err != nil {
		yakit.Error("恢复文件指针失败：%v", err)
		die(err)
	}
}


for fileReader.Next(){
    var currentTarget  = fileReader.Text()
	yakit.StatusCard("当前文件指针",fileReader.GetLastRecordPtr())
    swg.Add()
    go func {
        defer swg.Done()
		addScanHostCount()
		wg = sync.NewWaitGroup()
        wg.Add(1)
        go func {
            defer wg.Done()

            handleTCP(currentTarget)
        }
		if synscanEnable {
			wg.Add(1)
			go func {
				defer wg.Done()
				handleSYN(currentTarget)
        	}
		}
        wg.Wait()
    }
}



synWg.Wait()
swg.Wait()
yakit.SetProgress(0.94)

go fn{
	for {
		yakit.Info("漏洞检测中...")
		sleep(0.2)
		yakit.Info("漏洞检测中.....")
		sleep(0.2)
		yakit.Info("漏洞检测中........")
		sleep(0.1)
	}
}
manager.Wait()
yakit.SetProgress(1)
yakit.Info("本次扫描任务已完成。")
`

func (s *Server) SimpleDetect(req *ypb.RecordPortScanRequest, stream ypb.Yak_SimpleDetectServer) error {
	reqParams := &ypb.ExecRequest{
		Script: simpleDetect,
	}
	reqRecord := req.LastRecord
	reqPortScan := req.PortScanRequest

	// 把文件写到本地。
	tmpTargetFile, err := ioutil.TempFile("", "yakit-portscan-*.txt")
	if err != nil {
		return utils.Errorf("create temp target file failed: %s", err)
	}
	var targetsLineFromFile []string
	filePaths := utils.PrettifyListFromStringSplited(reqPortScan.GetTargetsFile(), ",")
	for _, filePath := range filePaths {
		raw, _ := ioutil.ReadFile(filePath)
		targetsLineFromFile = append(targetsLineFromFile, utils.PrettifyListFromStringSplited(string(raw), "\n")...)
	}

	targetsLine := utils.PrettifyListFromStringSplited(reqPortScan.GetTargets(), "\n")
	targets := append(targetsLine, targetsLineFromFile...)
	var allTargets string
	for _, target := range targets {
		hosts := utils.ParseStringToHosts(target)
		// 如果长度为 1 , 说明就是单个 IP
		if len(hosts) == 1 {
			allTargets += hosts[0] + "\n"
		} else {
			allTargets += strings.Join(hosts, "\n")
		}
	}
	allTargets = strings.Trim(allTargets, "\n")
	hostCount := len(strings.Split(allTargets, "\n"))
	if reqPortScan.GetEnableCClassScan() {
		allTargets = utils.ParseStringToCClassHosts(allTargets)
	}
	_, _ = tmpTargetFile.WriteString(allTargets)
	if len(targets) <= 0 {
		return utils.Errorf("empty targets")
	}
	tmpTargetFile.Close()
	defer os.RemoveAll(tmpTargetFile.Name())

	reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
		Key:   "target-file",
		Value: tmpTargetFile.Name(),
	})

	reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
		Key:   "task-name",
		Value: reqPortScan.TaskName,
	})

	reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
		Key:   "host-count",
		Value: fmt.Sprintf("%v", hostCount),
	})

	if reqRecord.Percent > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "record-ptr",
			Value: strconv.FormatInt(reqRecord.GetLastRecordPtr(), 10),
		})

		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "progress-data",
			Value: fmt.Sprintf("%.3f", reqRecord.GetPercent()),
		})
	}

	reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
		Key:   "ports",
		Value: reqPortScan.Ports,
	})
	reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
		Key:   "mode",
		Value: reqPortScan.GetMode(),
	})

	if reqPortScan.GetExcludeHosts() != "" {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "exclude-hosts",
			Value: reqPortScan.GetExcludeHosts(),
		})
	}

	if reqPortScan.GetExcludePorts() != "" {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "exclude-ports",
			Value: reqPortScan.GetExcludePorts(),
		})
	}

	if reqPortScan.GetSaveToDB() {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key: "save-to-db",
		})
	}

	if reqPortScan.GetSaveClosedPorts() {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key: "save-closed-ports",
		})
	}

	// 主动发包
	if reqPortScan.GetActive() {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key: "active",
		})
	}

	// 设置指纹扫描的并发
	if reqPortScan.GetConcurrent() > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "concurrent",
			Value: fmt.Sprint(reqPortScan.GetConcurrent()),
		})
	} else {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "concurrent",
			Value: fmt.Sprint(50),
		})
	}

	// 设置 SYN 扫描的并发
	if reqPortScan.GetSynConcurrent() > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "syn-concurrent", Value: fmt.Sprint(reqPortScan.GetSynConcurrent())})
	} else {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "syn-concurrent", Value: "1000"})
	}

	if len(reqPortScan.GetProto()) > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "proto",
			Value: strings.Join(reqPortScan.GetProto(), ","),
		})
	} else {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "proto",
			Value: "tcp",
		})
	}

	if len(utils.StringArrayFilterEmpty(reqPortScan.GetProxy())) > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "proxy",
			Value: strings.Join(reqPortScan.GetProxy(), ","),
		})
	}

	// 爬虫设置
	if reqPortScan.GetEnableBasicCrawler() {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key: "enable-basic-crawler",
		})
	}
	if reqPortScan.GetBasicCrawlerRequestMax() > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "basic-crawler-request-max",
			Value: fmt.Sprint(reqPortScan.GetBasicCrawlerRequestMax()),
		})
	} else {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "basic-crawler-request-max",
			Value: "5",
		})
	}

	if reqPortScan.GetProbeTimeout() > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "probe-timeout", Value: fmt.Sprint(reqPortScan.GetProbeTimeout())})
	} else {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "probe-timeout", Value: "5.0"})
	}

	if reqPortScan.GetProbeMax() > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "probe-max", Value: "3"})
	}

	switch reqPortScan.GetFingerprintMode() {
	case "service":
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "fp-mode",
			Value: "service",
		})
	case "web":
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "fp-mode",
			Value: "web",
		})
	default:
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{
			Key:   "fp-mode",
			Value: "all",
		})
	}

	// handle plugin names
	var callback func()
	reqParams.Params, callback, err = appendPluginNamesEx("script-name-file", "\n", reqParams.Params, reqPortScan.GetScriptNames()...)
	if callback != nil {
		defer callback()
	}
	if err != nil {
		return utils.Errorf("load plugin names failed: %s", err)
	}

	if reqPortScan.GetSkippedHostAliveScan() {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "skipped-host-alive-scan"})
	}

	if reqPortScan.GetHostAliveConcurrent() > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "host-alive-concurrent", Value: fmt.Sprint(reqPortScan.GetHostAliveConcurrent())})
	} else {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "host-alive-concurrent", Value: fmt.Sprint(20)})
	}

	if reqPortScan.GetHostAliveTimeout() > 0 {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "host-alive-timeout", Value: fmt.Sprint(reqPortScan.GetHostAliveTimeout())})
	} else {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "host-alive-timeout", Value: fmt.Sprint(5.0)})
	}

	if reqPortScan.GetHostAlivePorts() != "" {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "host-alive-ports", Value: fmt.Sprint(reqPortScan.GetHostAlivePorts())})
	} else {
		reqParams.Params = append(reqParams.Params, &ypb.ExecParamItem{Key: "host-alive-ports", Value: "22,80,443"})
	}

	return s.Exec(reqParams, stream)
}

func (s *Server) SaveCancelSimpleDetect(ctx context.Context, req *ypb.RecordPortScanRequest) (*ypb.Empty, error) {
	// 用于管理进度保存相关内容
	manager := NewProgressManager(s.GetProjectDatabase())
	uid := uuid.NewV4().String()
	manager.AddSimpleDetectTaskToPool(uid, req)
	return nil, nil
}

func (s *Server) RecoverSimpleDetectUnfinishedTask(req *ypb.RecoverExecBatchYakScriptUnfinishedTaskRequest, stream ypb.Yak_RecoverSimpleDetectUnfinishedTaskServer) error {
	manager := NewProgressManager(s.GetProjectDatabase())
	reqTask, err := manager.GetSimpleProgressByUid(req.GetUid(), true)
	if err != nil {
		return utils.Errorf("recover request by uid[%s] failed: %s", req.GetUid(), err)
	}

	return s.SimpleDetect(reqTask, stream)
}
