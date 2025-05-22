create-html:
	make create-html-config0
	make create-html-config1
	make create-html-config2
	cd perfomance-result/bins && vegeta plot \
      -title="Combined Throughput: Config0 vs Config1 vs Config2" \
      load_config0.bin \
      load_config1.bin \
      load_config2.bin \
      > ../html/combined_throughput.html
create-html-config0:
	cd perfomance-result/bins && vegeta plot \
          -title="Config0" \
          load_config0.bin \
          > ../html/config0.html
create-html-config1:
	cd perfomance-result/bins && vegeta plot \
              -title="Config1" \
              load_config1.bin \
              > ../html/config1.html
create-html-config2:
	cd perfomance-result/bins && vegeta plot \
              -title="Config2" \
              load_config2.bin \
              > ../html/config2.html
open-connection:
	ssh -p 2222 \
      -L 18081:stload.se.ifmo.ru:8080 \
      -L 18082:stload.se.ifmo.ru:8080 \
      -L 18083:stload.se.ifmo.ru:8080 \
      s368601@se.ifmo.ru