# Raspberry PI DS Cluster Notes (3/5/2025)

- Guide on how to set up a DS Cluster using Raspberry Pis

## MPICH

- MPI for high performance/quality computing

## Setup Steps

### Step 1

- connect up multiple Raspberry Pis via ethernet cable

### Step 2

- set up MPICH on the devices (thru Raspberry Pi OS or some other)
- update the nodes:

bashCopy codesudo apt-get update
sudo apt-get upgrade

- install mpich

bashCopy codesudo apt-get install mpich

- run some test script to ensure its working

### Step 3

- setup ssh keys
- setup machinefile
  - machinefile = text doc that details ip addresses of nodes
  - lets MPICH have visibility on full cluster

## Applications

- parallel programming
- distributed data processing
- automation
- high-computation simulations

## Optimizations

- central control node
- overclocking
