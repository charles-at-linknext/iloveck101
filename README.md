# iloveck101

致敬，[https://github.com/tzangms/iloveck101](https://github.com/tzangms/iloveck101), golang version

新增 Features: parallel downloading

iloveck101 我愛卡提諾
====================

沒錯, 你沒看錯, 我愛卡提諾。

我只想看圖, 可是卡提諾的網站載入太慢, 一堆廣告跟不必要的資訊, 還要等圖片 lazy load, 所以寫了這個 iloveck101, 直接先把圖片都下載下來, 看圖快多了。

目前只支援 Mac, Linux 我不確定。

by tzangms aka 海總理 2012.12.02



如何安裝
==========

```bash
$ go get github.com/lazywei/iloveck101
```


如何使用
===========

```bash
$ iloveck101 -u [url] -w [workers number]

```

for example, 你可以試試看思穎

```bash
$ iloveck101 -u http://ck101.com/thread-2876990-1-1.html -w 50
```

最後你可以在你 Mac 的「圖片」資料夾裡發現 `iloveck101` 的目錄, 你會很驚訝, 怎麼圖片都已經抓好放在裡面了呢?! 
沒錯, 這就是對卡提諾的愛啊!


看起來像這樣
==============

![圖片資料夾](https://raw.github.com/tzangms/iloveck101/master/docs/images1.png) 
![圖片](https://raw.github.com/tzangms/iloveck101/master/docs/images2.png) 


## License
Copyright (c) 2013 clonn
Licensed under the MIT license.
