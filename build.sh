#!/usr/bin/env bash

package="github.com/wassimbj/lnk"
package_name="lnk"
version=$1

if [[ -z version ]]; then 
   echo "Please give the version of the package"
   exit 1
fi

platforms=(
   "windows/amd64"
   "windows/386"
   "freebsd/amd64"
   "openbsd/amd64"
   "linux/amd64"
   "linux/arm64"
   "darwin/amd64"
   "darwin/arm64"
)

# build dir will contain the compressed executable build 
if [[ -n "build" ]]; then
   mkdir build
fi

for platform in "${platforms[@]}"
do
   platform_split=(${platform//\// })
   GOOS=${platform_split[0]}
   GOARCH=${platform_split[1]}

   output_name="$package_name-$version-$GOOS-$GOARCH"
   z_output_name="$package_name-$version-$GOOS-$GOARCH"

   if [[ $GOOS == "windows" ]]; then
      output_name+=".exe"
   fi

   env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package 

   if [[ $? != 0 ]]; then
      echo "Error, something went wrong"
      exit 1
   fi

   # zip folder/file ...
   echo "Building..."
   if [[ $GOOS == "linux" || $GOOS == "openbsd" ||  $GOOS == "freebsd" ]]; then
      tar -a -c -f "$z_output_name.tar" $output_name
      gzip "$z_output_name.tar"
      rm "$z_output_name.tar"
      if [[ -e "./build/$z_output_name.tar.gz" ]]; then
         echo "Remove existing build..."
         rm "./build/$z_output_name.tar.gz"
      fi
      echo "Moving to the build dir..."
      mv "$z_output_name.tar.gz" build
      rm $output_name
   else
      # use 7-Zip for windows else use zip
      if [[ $OSTYPE == "msys" ]]; then 
         7z a -tzip "$z_output_name.zip" $output_name
      else
         zip "$z_output_name.zip" $output_name
      fi
      if [[ -e "./build/$z_output_name.zip" ]]; then
         echo "Remove existing build..."
         rm "./build/$z_output_name.zip"
      fi
      echo "Moving to the build dir..."
      mv "$z_output_name.zip" build
      rm $output_name
   fi

   if [[ $? == 0 ]]; then
      echo "Success !!"
   else
      echo "Oops ! something went wrong"
   fi

   
done
