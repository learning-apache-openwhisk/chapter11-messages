IMG=actionloop/golang-v1.11:3

send.done: main.zip
	wsk action update mesg/send main.zip \
	--docker $(IMG) -P ../../cred.json && \
	touch send.done

main.zip: main_src.zip 
	docker run -i $(IMG) -compile main <$< >$@

main_src.zip: | main/vendor
	zip -r $@ main

main/vendor:
	cd main ; GOPATH=$(PWD)/../.. dep ensure

clean:
	rm *.zip *.done
