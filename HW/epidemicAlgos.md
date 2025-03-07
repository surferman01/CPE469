# Epidemic Algorithm Notes

Epidemic Algorithms for Replicated Database Management

## Problem

- Efficiently updating multiple copied databases with new information and ensuring they are synced

### Factors

- time required for update to propagate to all sites
- network traffic generated by single update

### Trials

1. Direct mail: send the updated contents directly to every other database
2. Anti-entropy: every site randomly chooses another and exchanged database data, resolving differences between the 2
3. Rumor mongering: site with new data periodically, randomly, chooses another site to share data with until it has tried the same spot x times

### Terminology

- Infective: site with new data that it wants to share
- Susceptable: site that does not have new info
- Removed: recieved update and is not willing to share

- anti-entropy is a site with only infective or susceptable clients
ex)
When anti-entropy is used as a backup mechanism, there is a significant advantage in using a complex epidemic such as rumor mongering rather than direct mail for initial distribution of updates.

- residue: how much longer after everyone has gotten the update does data linger / checks still occur
- traffic: how much network traffic impact the protocol creates

ex)
Recall that with respect to an individual update, a database is either susceptible (it does not know the update), infective (it knows the update and is actively sharing it with others), or removed (it knows the update but is not spreading it)

## Push - Pull

- Like github, determining whether to take data from one source or another based on the timestamp
- use pull or push-pull rather than push when anti-entropy is used as a backup to some other distribution mechanism
  - this is very taxing
  - use checksums so we dont need to constantly send the entire db over the network

ex)
Two sites s and s' perform anti-entropy by first exchanging
recent update lists, using the lists to update their databases and checksums, and then comparing checksums. Only if the checksums disagree do the sites compare their entire databases

- for this, use a recent update list so that the changes which is bound by the expected distribution time

## Gossip / Rumor spreading

### Variations

- feedback vs blind
  - feedback
    - sender loses interest in the rumor if the recipient does not respond to it (it already has heard the rumor)
  - Blind
    - There is a 1/k chance that the rumor is dropped by sender regardless of response

- Counter vs coin
  - counter
    - instead of 1/k chance, make it lose interest if k recipients dont respond
    - combine counter with blind to make it infective regardless of responses

- Push vs pull
  - pull is good since it will find information that has multiple rumors heard

### Deletion / Death Certs

- When you remove an entry, you need to propagate it as though it was data since table removals will be thought of as out of date entries
  - give it a 'death certificate' w a timestamp
  - now how do we re-use this space?
    - one way is to just wait for the cert. to replace all the older data on other machines and then let it expire after x time (days perhaps)

- give death certs. an activiation time as well to prevent accidentally reinstating a deleted item since adding a death cert to a new database updates its original timestamp

## Entry for canvas

This article provides thorough insight and examples into the application of epidemic algorithms for propagating data through air-gapped databases. The author begins by discussing the problems faced: how to efficiently update multiple (large) databases that are expected to share information identically with the least network impact and 'residue'. They then begin to introduce the primary factors to be considered such as network traffic, time to propagate data to the final database, and the 'residue'. To understand these terms, one needs to understand the different means of propagating data.

1. Direct mail: send the updated contents directly to every other database

    - problems: this requires large network transfers as the entire database needs to be passed over the network. The database with the update may not have visibility on all databases that require updating.
    - pros: deterministic outcome since we know where to send and how much to send to each database

2. Anti-entropy: every site randomly chooses another and exchanges database data, resolving differences between the 2

    - problems: we still need to pass the entire database across the network. slower; since randomness is involved, it may take longer to get to the last database
    - pros: get visibility on all the databases since the data will continue to propagate (only infectors and susceptible)

3. Rumor mongering: A site with new data periodically, randomly, chooses another site to share data with until it has tried the same spot x times. The data will become 'hot' and passed around to random databases until it has been 'hot' for too long (the sender doesn't receive x responses)

    - problems: this may take a long time since randomness exists.
    - pros: the rumor can use a recent updates table that allows minimal network traffic when a rumor is passed and the databases are updated.

Some more terminology that is important is as follows:

- Infective: site with new data that it wants to share
- Susceptible: site that does not have new info
- Removed: received update and is not willing to share

These terms relate to the status of the database related to the updated data. Typically, an infective database is one that will be sending gossip or data to another. The susceptible is a database that does not have the new information. Removed is the state for a database where the rumor has stopped spreading.

While the summary thus far examines the different methods of data propagation, the article also delves into optimizing the network traffic and residue of the algorithms by changing parameters such as how many requests it takes for a rumor to go away (residue minimization) and using update tables (network bw minimization).

The final important part that the article investigates is the removal of entries from a database and the propagation of that data. To do this, we use a 'death certificate' that is similar to another data input, however, it has 2-time parameters: the time it was added and the activation time. These parameters are used to ensure we don't re-add previously deleted data (that may still exist on a database), and to determine when the space will free up and delete itself (typically some x amount of days). If we knew exactly how long it would take for the death certificate to get to each database, then we would simply delete the entry at that exact time– since we don't know, we apply some expected time.
