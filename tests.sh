echo "########### CLEARING GO CACHE ###############"
go clean -cache

echo "########### STARTING TEST SUITE #############"
ginkgo -r -focus "gostorz"
ginkgo -r -focus "mgen"

ginkgo -r -focus "cache"
ginkgo -r -focus "react"
ginkgo -r -focus "client"

cd test
./tests.sh
cd ..
