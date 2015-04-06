#!/bin/bash
URL="http://mangafox.me/manga"
DIR="manga"
################## Function #####################
cnt_words()
{

        file=$1
        srch="$2"
        v=$(( (`cat $file | wc -c` - `sed "s/$srch//g" $file | wc -c`) / ${#srch} ))
        echo $v
}

find_max_page()
{
	cat 1page.html |grep option|head -1|awk -F"<option value=\"0\"" '{print $1}'>2.html
	MAXPAGE=`cnt_words 2.html option`
	MAXPAGE=`expr $MAXPAGE / 2`
	rm 2.html
	echo $MAXPAGE
}

find_max_chapter()
{
	VOL=`enc_vol $1`
	MAX_CHAP=`cat 0page.html |grep $VOL|wc -l|awk -F" " '{print $1}'`
	echo $MAX_CHAP
}

find_max_volume()
{
	MAX_VOL=`cat 0page.html |grep volume |wc -l|awk -F" " '{print $1}'`
	echo $MAX_VOL
}

enc_vol()
{
	if [ $1 -le 9 ]
	then
		echo "v0$1"
	else
		echo "v$1"
	fi
}

enc_chap()
{
	CHAP=`cat 0page.html|grep $1|awk -F"1.html" '{print $1}'|awk -F"$1" '{print $2}'|awk -F"/" '{print $2}'|tail -$2|head -1`
	echo $CHAP
}

load()
{
#echo "$1 $2 $3 $4 $5 $URL"
#     name cap page start  end
#     name $CHAP $NOVOL 1 $MAXPAGE
if [ "$3" == "NOVOLUME" ]
then
	mkdir -p $DIR/$1/$3/$2
else
	mkdir -p $DIR/$1/$2/$3
fi
for ((k=$4 ; k<= $5 ;k++))
do
	if [ "$3" == "NOVOLUME" ]
	then
		curl -s "$URL/$1/$2/$k.html" > 1page.html
	else
		curl -s "$URL/$1/$2/$3/$k.html" > 1page.html
	fi
	cat 1page.html |grep "mfcdn.net/store" |grep -v thumbnails |awk -F"src=\"" '{print $2}'|awk -F"\"" '{print $1}'|head -1 >data.txt
	MANGA=`cat data.txt`
	echo "Download[$k/$5]: $MANGA"
	if [ "$3" == "NOVOLUME" ]
	then
		wget --quiet -P $DIR/$1/$3/$2 $MANGA
	else
		wget --quiet -P $DIR/$1/$2/$3 $MANGA
	fi
	rm data.txt 1page.html
done
}

################## Main #####################
if [ $# -eq 0 ]
then
	echo "Use: $0 [Serie] [Volume] [Chapter] [PageStart] [PageEnd]"
	echo "Example: ./mangafox.sh btooom v01 c001 20 60"
	echo "Example: ./mangafox.sh btooom c001 c017 17"
	echo "Example: ./mangafox.sh btooom v01 c001"
	echo "Example: ./mangafox.sh btooom v09"
	echo "Example: ./mangafox.sh btooom"
	exit
fi

if [ $# -eq 1 ]
then
	echo -n "Finding max volume ... "
	curl -s "$URL/$1/" > 0page.html
	MAX_VOL=`find_max_volume`
	echo "$MAX_VOL"
	for ((i=1 ; i<= $MAX_VOL ; i++))
	do
		echo -n "Finding max chapter of Volume[$i] ... "
		MAX_CHAP=`find_max_chapter $i`
		echo "$MAX_CHAP"
		for ((j=1 ; j<= $MAX_CHAP ; j++))
		do
			echo "Loading from [Volume..$i] [Chapter..$j]"
			VOL=`enc_vol $i`
			CHAP=`enc_chap $VOL $j`
			#echo "Get Page of $URL/$1/$VOL${CHAP}1.html"
			curl -s "$URL/$1/$VOL/$CHAP/1.html" > 1page.html
			echo -n "Finding max page of [Volume $i] [Chapter $j] ... "
			MAXPAGE=`find_max_page`
			echo $MAXPAGE
			load $1 $VOL $CHAP 1 $MAXPAGE
		done
	done
	exit
fi

if [ $# -eq 2 ]
then
	curl -s "$URL/$1/" > 0page.html
	echo -n "Finding max chapter of Volume[$2] ... "
	MAX_CHAP=`cat 0page.html |grep $2|wc -l|awk -F" " '{print $1}'`
	#MAX_CHAP=`find_max_chapter $2`
	echo "$MAX_CHAP"
	for ((j=1 ; j<= $MAX_CHAP ; j++))
	do
		echo "Loading from [Volume..$2] [Chapter..$j]"
		VOL=`echo $2`
		CHAP=`enc_chap $VOL $j`
		#echo "Get Page of $URL/$1/$VOL${CHAP}1.html"
		curl -s "$URL/$1/$VOL/$CHAP/1.html" > 1page.html
		echo -n "Finding max page of [Volume $2] [Chapter $j] ... "
		MAXPAGE=`find_max_page`
		echo $MAXPAGE
		load $1 $VOL $CHAP 1 $MAXPAGE
	done
	
fi

if [ $# -eq 3 ]
then
	curl -s "$URL/$1/$2/$3/1.html" > 1page.html
	echo -n "Finding max page ... "
	MAXPAGE=`find_max_page`
	echo $MAXPAGE
	load $1 $2 $3 1 $MAXPAGE
	exit
fi

if [ $# -eq 4 ]
then
	#./mangafox.sh btooom c001 c017 17"
	#$0		$1	$2  $3  $4
	for ((j=1 ; j<= $4 ; j++))
	do
		echo "Loading from [Chapter ..$j]"
		if [ $j -le 9 ]
		then
			CHAP="c00$j"
		elif [ $j -le 99 ]
		then
			CHAP="c0$j"
		else
			CHAP="c$j"
		fi
		curl -s "$URL/$1/$CHAP/1.html" > 1page.html
		echo -n "Finding max page ..."
		MAXPAGE=`find_max_page`
		echo $MAXPAGE
		NOVOL="NOVOLUME"
		load $1 $CHAP $NOVOL 1 $MAXPAGE
	done
fi

if [ $# -eq 5 ]
then
	load $1 $2 $3 $4 $5
fi

################## End #####################
	
