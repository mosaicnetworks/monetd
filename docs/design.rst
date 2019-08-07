.. _design_rst:

Design
======

In this document we explain our implementation of the MONET Hub; in particular
the mechanism that dictates who can participate in the consensus system, and
how to make participants accountable for their actions. Before deliberating on
an implementation, it is important to have a clear picture of the desired
outcome. So we will start by reiterating the role of the Hub in MONET, and
outline its principal requirements. We then visit the spectrum of potential
implementations before explaining our choice of a permissioned Byzantine Fault
Tolerant (**BFT**) consensus algorithm coupled to the Ethereum Virtual Machine
(**EVM**). Lastly we weigh up the pros and cons of Proof of Stake (**PoS**),
and explain our decision to implement PoA for the time being.

MONET and the MONET Hub
-----------------------

MONET’s mission is to boost the adoption of peer-to-peer architectures by
enabling mobile devices to connect directly to one another in dynamic ad-hoc
networks. We believe that a new generation of applications will emerge from
this technology. The real force behind MONET, which makes it original and
disruptive, is the concept of Mobile Ad-Hoc Blockchains, and the open-source
software which implements it; particularly Babble, the powerful consensus
algorithm which is suitable for mobile deployments due to its speed,
bandwidth efficiency, and leaderlessness.

We anticipate that many MONET applications will require a common set of
services to persist non-transient data, carry information across ad-hoc
blockchains, and facilitate peer-discovery. So we set out to build the
MONET Hub, an additional public utility that provides these services. In the
spirit of open architecture, MONET doesn’t rely on any central authority,
so anyone is free to implement their own alternative, but the MONET Hub is
there to offer a reliable, fast, and secure solution to kickstart the system.

As such, the qualitative requirements of the Hub are:

+ **Speed**: It should support thousands of commands per second, with latencies
  under one	second.

+ **Finality**: Results from the hub should be definitive, without the
  possibility of being arbitrarily overridden in the future.

+ **Availability**: It should provide a continuous service in the face of
  network failures or isolated disruptions.

+ **Cost**: As we want to lower the barrier to entry for developers, using the
  Hub should be cheaper than rolling out one’s own solution.

+ **Security**: The hub should provide a trusted source of data and computation,
  with measures guarding against information loss, data manipulation, or
  censorship.

+ **Governance**: The set of entities controlling this utility should be
  transparent, with a mechanism to add or remove participants, and keep them
  accountable for their actions.

+ **Flexibility**: It should be possible and relatively easy to update the
  software, recover from failures, and adapt to changes.

Spectrum of possible Implementations
------------------------------------

From a simple web-service hosted on a privately-owned server, to a public
global blockchain like Ethereum, there are many potential ways to implement
this service. However, given our requirements, a simple server scores pretty
low in all categories (except perhaps speed and flexibility), and global public
blockchains are too slow, too hard to update, and usually provide only
probabilistic finality, which is not acceptable.

Somewhere in the middle lies a category of distributed systems consisting of
relatively small clusters of servers maintaining identical copies of an
application via sophisticated communication routines and consensus algorithms.
Within this category, there are instances where the entire cluster is
controlled by a single entity, and others where each replica is controlled
by a different entity.

Modern blockchain projects, including cryptocurrencies like Facebook’s Libra
and the Cosmos Atom, adopt the second variant, where nodes are controlled by
different entities. A naive implementation would render them vulnerable to
malicious actors trying to subvert the system; hence they require strong
consensus algorithms, commonly referred to as Byzantine Fault Tolerant (BFT),
and a reputation system to incentivize good behavior and punish malicious
actors.

Given the requirements stated in the previous section, we believe that the
MONET Hub falls in the same category, and requires a permissioned BFT system.

Ethereum with Babble Consensus
------------------------------

We have developed the Monet Toolchain, a complete set of software tools for
setting up and using the MONET Hub. This includes ``monetd``, the software
daemon that powers nodes on the MONET Hub.

To build ``monetd``, we used our own BFT consensus algorithm, `Babble
<https://github.com/mosaicnetworks/babble>`__, because it is fast, leaderless,
and offers finality. For the application state and smart-contract platform, we
use the Ethereum Virtual Mahcine (EVM) via `EVM-Lite
<https://github.com/mosaicnetworks/evm-lite>`__, which is a stripped down
version of `Go-Ethereum <https://github.com/ethereum/go-ethereum>`__.

The EVM is a security-oriented virtual machine specifically designed to run
untrusted code on a network of computers. Every transaction applied to the EVM
modifies the State which is persisted in a Merkle Patricia tree. This data
structure allows to simply check if a given transaction was actually applied to
the VM and can reduce the entire State to a single hash (merkle root) rather
analogous to a fingerprint.

The EVM is meant to be used in conjunction with a system that broadcasts
transactions across network participants and ensures that everyone executes the
same transactions in the same order. Ethereum uses a Blockchain and a Proof of
Work consensus algorithm. EVM-Lite makes it easy to use any consensus system,
including `Babble <https://github.com/mosaicnetworks/babble>`__.

The remaining question is how to govern the validator-set, and what to use as a
reputation system to punish or incentivise participants to behave correctly.

PoS and PoA
-----------

A BFT consensus algorithm ensures that a distributed system remains available
and consistent in adversarial conditions, with some nodes exhibiting arbitrary
failures or malicious behavior, as long as a majority of participants are
functioning correctly (actually ⅔). Any trust in the system therefore depends
on the ability to legitimise this assumption. What is needed is a mechanism to
ensure, with a high degree of confidence, that at least two thirds of
participants in the consensus system are functioning correctly at all times.
The problem is two-fold: who gets to be a participant, and how are participants
incentivised to behave correctly? Not surprisingly, the most convincing answers
revolve around money or reputational risk.

In a Proof of Stake (PoS) arrangement, participants are required to lock a
significant portion of their assets (usually the blockchain’s built-in token),
and respect an extended un-bonding period when they want to leave. At any given
time, the validator set is defined by the top N stakers, where N is the desired
size of the validator-set. If they are caught undermining the network, this
deposit is destroyed. Hence, participants are deterred from cheating.
Additionally, participants are usually programmatically compensated for
actively participating in securing the network. Hence they are incentivised to
act correctly. A nice feature of PoS is that, being a very capitalistic model,
it is relatively open; anyone can participate without asking for permission,
as long as they put up a stake.

In Proof of Authority (PoA), the stake is tied to reputational risk. It relies
on the natural aversion of most humans to tarnish their own reputation.
The list of allowed validators is governed by a whitelist. The whitelist is
amended through a voting process among existing whitelisted entities. This
scheme is less anonymous or open than PoS but has deep roots. The trust of a
PoA system rests on the initial group of participants because any amendment
to the list has to gather consensus from them;
so the trust (or distrust) is carried over as the
validator-set evolves. In a system like Babble, the most serious offence
consists in signing two different blocks at the same height. Evidence of this
can be packaged into an irrefutable proof, and used to punish the guilty
participants.

Proof of Stake opens exciting opportunities for a variety of stakeholders, and
these economic incentives are excellent for the industry as they drive
innovation. That being said, we are of the opinion that it is too early to
ascertain the resilience of PoS in the face of decisive attacks, as current
production deployments are very recent, and the theoretical arguments alone are
not sufficiently convincing (although they sound quite reasonable). We are
keeping an eye on PoS systems, hoping that they withstand the test of time. In
the meantime, we have opted to implement PoA, to roll out a reliable version of
the MONET Hub, with an eye on extending to PoS in a coordinated software update
later down the road.

Conclusion
----------

The MONET Hub is a pivotal utility that facilitates the creation of mobile
ad-hoc blockchains, and the emergence of a new breed of decentralised
applications. To maximise the performance, security, and flexibility of this
system, we have opted to build the Monet Toolchain, a smart-contract platform based on the Ethereum
Virtual Machine and a state-of-the-art BFT consensus algorithm, Babble. To
govern the validator-set involved in the consensus algorithm, we have chosen to
implement a Proof of Authority system, with the idea of extending to Proof of
Stake when more evidence of its efficacy becomes available.
