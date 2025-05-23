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

create-histograms:
	mkdir -p perfomance-result/histograms
	cd perfomance-result && \
	for f in bins/load_config*.bin; do \
	  base=$$(basename $$f .bin); \
	  vegeta report \
		-type=hist\[0,100ms,200ms,300ms,400ms,500ms,600ms,670ms,1s\] \
		$$f > histograms/$${base}_latency_hist.txt; \
	done
create-reports:
	mkdir -p perfomance-result/reports
	cd perfomance-result && \
	for f in bins/load_config*.bin; do \
	  base=$$(basename $$f .bin); \
	  vegeta report -type=json $$f | jq . > reports/$${base}.json; \
	done
create-csv:
	echo "config,mean_req_per_s" > perfomance-result/mean_throughput.csv

	cd perfomance-result && \
	for f in reports/*.json; do \
	  base=$$(basename $$f .json); \
	  mean=$$(jq .throughput $$f); \
	  echo "$${base},$${mean}" >> mean_throughput.csv; \
	done


show:
	make create-html
	make create-histograms
	make create-reports
	make create-csv
open-connection:
	ssh -p 2222 \
      -L 18081:stload.se.ifmo.ru:8080 \
      -L 18082:stload.se.ifmo.ru:8080 \
      -L 18083:stload.se.ifmo.ru:8080 \
      s368601@se.ifmo.ru