#!/bin/bash
if [ -z $2 ]
then
	echo "use $0 [manga name] [maxchapter]"
	echo "example $0 read-inu-yashiki-chapter 28"
else
	echo "Build dir for save file $1"
	for ((volume=1 ; volume<=$2 ;volume++))
	do
		echo "Build dir for save file $1/$volume"
		mkdir -p $1/$volume
		echo "Loading page... [$volume/$2]"
		wget --quiet http://www.thai-cartoon.net/$1-$volume.html -O volume.html
		cat volume.html |grep "</Script>'"|awk -F"'" '{print $2}'>img.txt
		COUNT=1
		cat img.txt|while read LINE
		do
			echo "Loading... image $COUNT.jpg in [$volume/$2]"
			wget --quiet $LINE -O $1/$volume/$COUNT.jpg
			COUNT=`expr $COUNT + 1`
		done
	done
	rm volume.html
	rm img.txt
fi
