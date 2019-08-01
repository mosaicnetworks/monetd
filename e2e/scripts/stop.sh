#!/bin/bash

NET=${1:-"net9"}

giverny network stop --remove $NET 