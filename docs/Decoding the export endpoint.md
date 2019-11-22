# Decoding the Export Endpoint

The ``/export`` endpoint generates a ``genesis.json`` format file. Within that
format is a ``Poa`` key - which holds details of the Proof of Acceptance Smart
Contract. Within the ``Poa`` key is a ``storage`` key which holds the EVM storage
for the smart contract. This document helps you decipher that. 

We are going to work through this sample from the testnet. Only part of the Poa key
is included here - the ABI and bytecode are rather large.
	
```json

	{
	"Poa": {
    "Address": "0xaBBAABbaaBbAABbaABbAABbAABbaAbbaaBbaaBBa",
    "Balance": "0",
    "Storage": {
	"0000000000000000000000000000000000000000000000000000000000000001": "06",
	"0000000000000000000000000000000000000000000000000000000000000002": "06",
	"0000000000000000000000000000000000000000000000000000000000000004": "03",
	"008bd263fb326b18ca849a3eaa2e862e4b38733eae50a1c136ca3d45f70d7f71": "9501bca3ec820659ff257c0cc134bce65b2c429017d9",
	"075915bd8cc73a3ec2c9e1f1cf5d1cf6a71b0a3a84ffe96af892b6adff606bc8": "a06d61777300000000000000000000000000000000000000000000000000000000",
	"07a406100bd1787d918d9b159c4e6af27da22009f8d8d32bc281449b54a5d722": "94dcbe3c8a1e9877214132ff167e99abecb6dd2ccc",
	"10c97a3499e0eed2046db286ab8e0a8084af62314b443f2fe0e00a5f5584299e": "94bca3ec820659ff257c0cc134bce65b2c429017d9",
	"1635eaee71f81fd82037f0e37a8d489f8435624f39d8f3bf1f37cb485396b299": "a04d61737465724a6f6e0000000000000000000000000000000000000000000000",
	"2380a9303b21b715f84eb99203b0053ab25aa7e6865d794ae19e718f2de1bc94": "9501da37bd41430c3242423d155bbd9b0e535067bbb3",
	"2450285b8dd318712b58ca998a50fa73c7b51caf3ae42ac83c0dc39cdc1f4f82": "950151e5e40e48b2a639f111e0867464d02b2f9f0325",
	"28999a6f81c03a0b3dc5146a7b4584c2d817b18eaa146e369faa58f876a1ae72": "9501dcbe3c8a1e9877214132ff167e99abecb6dd2ccc",
	"318fd9530f763482df6c5095c7dca17d04653dfb39f026b1b3340ba83f3d6209": "9501da37bd41430c3242423d155bbd9b0e535067bbb3",
	"3fd5a574a63b0ac663b143343d51af4da25a9ce8d9157196fe1b7bb03d8de1ae": "94da37bd41430c3242423d155bbd9b0e535067bbb3",
	"405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace": "94da37bd41430c3242423d155bbd9b0e535067bbb3",
	"405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf": "94bca3ec820659ff257c0cc134bce65b2c429017d9",
	"405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad0": "94640314f39674b19806de47bc67679e277196e9d4",
	"405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad1": "9451e5e40e48b2a639f111e0867464d02b2f9f0325",
	"405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad2": "94e7970cbd5912d7908c474ac6ae88000bb0c7cdf2",
	"405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad3": "94dcbe3c8a1e9877214132ff167e99abecb6dd2ccc",
	"41ff9ae19119a16933736c9bdc442385b911236cb67a7b2e75135cde84c71da9": "a06b6576696e000000000000000000000000000000000000000000000000000000",
	"54e08178835925742330520ecd07123463a88c3d5d875ed2e9d8700496ad6128": "9501640314f39674b19806de47bc67679e277196e9d4",
	"617c496a7750bd3858c74cfd95409f691062a00c76d1f9fa7f0d93747f9939da": "a070686f6e79320000000000000000000000000000000000000000000000000000",
	"716eb2c3e058a7aaebe4791d8d531e738f4d02b8f0081680d502955e309c12d6": "94bca3ec820659ff257c0cc134bce65b2c429017d9",
	"76a3f2f1b02499ab21e754e5afbe82347498ca07bb5b32965afc14d59c4348bb": "94da37bd41430c3242423d155bbd9b0e535067bbb3",
	"76a3f2f1b02499ab21e754e5afbe82347498ca07bb5b32965afc14d59c4348bc": "94bca3ec820659ff257c0cc134bce65b2c429017d9",
	"76a3f2f1b02499ab21e754e5afbe82347498ca07bb5b32965afc14d59c4348bd": "94640314f39674b19806de47bc67679e277196e9d4",
	"76a3f2f1b02499ab21e754e5afbe82347498ca07bb5b32965afc14d59c4348be": "94dcbe3c8a1e9877214132ff167e99abecb6dd2ccc",
	"76a3f2f1b02499ab21e754e5afbe82347498ca07bb5b32965afc14d59c4348bf": "9451e5e40e48b2a639f111e0867464d02b2f9f0325",
	"8200f7539727d20df27257fc84b577803d0c262709d841c214e35b6bf127d082": "94e7970cbd5912d7908c474ac6ae88000bb0c7cdf2",
	"86db9e829b3d4ea996dfcd2305737dc64a943e273074456517818c7ff13a6756": "9451e5e40e48b2a639f111e0867464d02b2f9f0325",
	"8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b": "94fde27dace5cd5b29a368aff83327c37c5ecb73da",
	"8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19c": "9485d5357a15a0ac954d3bd4b4770887af80f4d60b",
	"8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19d": "94b7f0d947278c3a49db3da85214716a61b5ca41ee",
	"ab2ac9c676326e06e2c4ab0ee623ca6b6c687fba894288525dd9f73ecf05d0e7": "9501da37bd41430c3242423d155bbd9b0e535067bbb3",
	"b1e77e65957c41e3e35620d293709b74c86c5984058214a97cfa6b49b07f7d42": "a064616e7500000000000000000000000000000000000000000000000000000000",
	"b78468f43e89a8a26c353d4054b0e2db604a0ad3e922607f9cbfbddeadc5a9c1": "a06d6f736169630000000000000000000000000000000000000000000000000000",
	"c32e8dbbed7463d206e12b4ff7ecd75520cd3e69494e508e2b52e3ef1753ac0b": "9501640314f39674b19806de47bc67679e277196e9d4",
	"ca9933d856026298445ae96b5f20a9e936470d7a378151227b47561edea0692b": "94640314f39674b19806de47bc67679e277196e9d4",
	"cb5df5bddb484c1c36d2aa99eed73eed1652b2ac5a0e9064d8d0f99bd352a38c": "94b7f0d947278c3a49db3da85214716a61b5ca41ee",
	"cb5df5bddb484c1c36d2aa99eed73eed1652b2ac5a0e9064d8d0f99bd352a38d": "94da37bd41430c3242423d155bbd9b0e535067bbb3",
	"cb5df5bddb484c1c36d2aa99eed73eed1652b2ac5a0e9064d8d0f99bd352a38e": "05",
	"cb5df5bddb484c1c36d2aa99eed73eed1652b2ac5a0e9064d8d0f99bd352a391": "05",
	"d103b9f532050d5f09d463c046b7190db9647dfb453a8439ffa535bb1c1f3ed9": "a06e616d6900000000000000000000000000000000000000000000000000000000",
	"dcb94dab15edc7d7f186f3d7347ef0d090c624b318143a92d9943f2c4aeda5ef": "a06d617274696e5f61727269766574730000000000000000000000000000000000",
	"ef36fef6435b8f26e4aa275cba21e75707a0781c3d8c7eb811db7f828e895146": "94640314f39674b19806de47bc67679e277196e9d4",
	"f7367c36e64c1e0c00d771a30b5d1c5a31c61c3e5b33ccb9f71e029088963882": "a0544b534c32000000000000000000000000000000000000000000000000000000",
	"fa2a575d669eac5bf10ed1013ebf47e633c95072080f9775a3fb534d8eed15ea": "94bca3ec820659ff257c0cc134bce65b2c429017d9"
    }
 }
 ```
    
The data is in key value pairs. The key is a slot, the value is a value. The slot where data is stored is deterministic. 
The same structure is always stored in the same place. 

Precisely where it is stored is covered well in [this article](https://programtheblockchain.com/posts/2018/03/09/understanding-ethereum-smart-contract-storage/?source=post_page-----d3383360ea1b----------------------)

The storage variables defined in the smart contract are:

```solidity
        mapping (address => WhitelistPerson) whiteList;
        uint whiteListCount;
        address[] whiteListArray;
        mapping (address => NomineeElection) nomineeList;
        address[] nomineeArray;
        mapping (address => bytes32) monikerList;
        mapping (address => NomineeElection) evictionList;
        address[] evictionArray;
```

Interpreting that list we get:

+ **whiteList** -- *slot 0* -- mapping, so no data stored
+ **whiteListCount** -- *slot 1* -- integer -- just the value is stored
+ **whiteListArray** -- *slot 2* -- array -- so length is stored
+ **nomineeList** -- *slot 3* -- mapping, so no data stored
+ **nomineeArray** -- *slot 4* -- array -- so length is stored
+ **monikerList** -- *slot 5* -- mapping, so no data stored
+ **evictionList** -- *slot 6* -- mapping, so no data stored
+ **evictionArray** -- *slot 7* -- array -- so length is stored

So at a crude level we would expect data in 1,2,4 and 7 -- with the other slots
empty because they are mappings, but we got (without the leading zeros):

```
	"1": "06",
	"2": "06",
	"4": "03",
```	

The missing "7" is because there are no entries on the eviction list, and zero
values are not stored in EVM. For any zero value, the relevant key will simply
be absent. So in this case, we have 6 on the whitelist and 3 on the nominee list.

There is a command ``giverny parse [genesis.json]`` that extracts the whitelist
entries from a genesis file. 

For this data, it yields:

```bash
$ giverny parse ~/Documents/export.json

POA Address:  0x0xaBBAABbaaBbAABbaABbAABbAABbaAbbaaBbaaBBa 

6 peers found 

0x640314f39674b19806de47bc67679e277196e9d4  danu
0x51e5e40e48b2a639f111e0867464d02b2f9f0325  MasterJon
0xdcbe3c8a1e9877214132ff167e99abecb6dd2ccc  nami
0xda37bd41430c3242423d155bbd9b0e535067bbb3  mosaic
0xe7970cbd5912d7908c474ac6ae88000bb0c7cdf2  kevin
0xbca3ec820659ff257c0cc134bce65b2c429017d9  martin_arrivets

Contract does not match the POA bytecode
This may not be an issue if a different release of Monetd was used to generate the genesis.json file.
Your version of Monetd is:
Monetd Version: 0.3.3-develop-93a1e090
     EVM-Lite Version: 0.3.6-develop
     Babble Version: 0.5.10-develop
     Geth Version: 1.8.27
Solc: Version: 0.5.11+commit.22be8592.Linux.g++ 
      Distributor ID:	Ubuntu;  	Description:	Ubuntu 18.04.3 LTS;  	Release:	18.04;  	Codename:	bionic;  
      develop ef45731 Development Version
```

It achieves this by looking at the keys from ``405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace``
onwards, reading the number of records denoted by the value in slot 2. The ``405...`` value is ``keccak256(2)`` 

To save you having to calculate them the values for the other arrays are:

+ Slot 2: ``0x405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace``
+ Slot 4: ``0x8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b``
+ Slot 7: ``0xa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688``

So if we look in slot ``0x8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b``
we find 3 consecutive items - the nominee list. 

You will note that the addresses have 42 or 44 characters - they are prefixed with
either ``94`` or ``9501``. You can convert to a regular address by taking the
rightmost 40 characters. 

Also note the absence of a record at ``0xa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688``
which is as we would expect -- there are no evictees. 

At this point the explainations got too involved. Take a look at the code in 
``cmd/giverny/commands/parse/parse.go`` if you want more detail. Otherwise
just run ``giverny parse [export output]``. The data section has the following 
columns:

+  Slot
+  Raw Data
+  Processed Data, i.e. addresses have prefixed removed, monikers are rendered as text
+  Description -- a narrative description of the field
+  Explained -- boolean flag as to whether the field is explained.


