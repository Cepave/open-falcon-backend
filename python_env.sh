#!/bin/bash

cd dashboard
virtualenv ./env --python=python2.7
source env/bin/activate
pip install -r pip_requirements.txt -i http://pypi.douban.com/simple --trusted-host pypi.douban.com

cd ../portal
virtualenv ./env --python=python2.7
source env/bin/activate
pip install -r pip_requirements.txt -i http://pypi.douban.com/simple --trusted-host pypi.douban.com