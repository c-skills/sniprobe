
<p align="center">
<img src="https://github.com/c-skills/welcome/blob/master/logo.jpg" />
</p>


sniprobe
========

This is not a political discussion. You will find censorship in almost any country
and it has many shapes. We have obviously chosen Iran in this example, for a good
and actual reason. Given the level of blocking in Iran as of today, SNI probing
might not be a solution anymore, but might be of help elsewhere. If you have lists
that worked for you or can do tests behind blocking equipement, let us now.

If you are interested in that topic and think censorship can only happen to others,
we recommend a look into the mirror and to read Bert's excellent [blog post](https://berthub.eu/articles/posts/who-controls-the-internet).


Intro
-----

`c->skills research` has conducted some research about censorship blockings in certain
countries. Our bros and anarchic noses over at [THC](https://blog.thc.org/the-iran-firewall-a-preliminary-report) have done so previously, too.
Despite that our own observations differ a bit (what is possible on day N might be impossible on day N+1),
they found out that at least in certain times the blocking is based on the SNI of the TLS
session that is seen by the censor.

According to our own tests, all major news sites ignore unknwon SNIs, so there is a chance
to abuse SNI based blocking to make them think that you are visiting a *legit* site,
while actually browsing behind the wall.

Probing
-------

First of all you want to build the package. It requires Go and C++ compilers to be installed,
just do `make`.

```
~ $ ./sniprobe www.bbc.co.uk 443 sni.txt
Success: https://www.bbc.co.uk:443 with SNI irangov.ir (Thu, 03 Nov 2022 08:40:00 GMT)
Success: https://www.bbc.co.uk:443 with SNI irna.ir (Thu, 03 Nov 2022 08:40:00 GMT)
Success: https://www.bbc.co.uk:443 with SNI en.irna.ir (Thu, 03 Nov 2022 08:40:00 GMT)
Success: https://www.bbc.co.uk:443 with SNI iranpress.com (Thu, 03 Nov 2022 08:40:00 GMT)
~ $
```

Ok, that works. If censored, you will see some error messages which are explained inside
`sni.txt`. Depending on who is blocking you, you have to build your own list of SNIs that
make sense and where you expect that these are white-listed by censor.

If you found a SNI that tells you `Success`, you have to edit `Makefile` and
set `GOOD_SNI` and the SSL library path that matches your system.

Then `make clean && make`.

sniswitcher
-----------

Now the spoiler. Large browsers (notably chrome) are nailed down to use their builtin TLS
implementations which is not easy to tamper with to set good SNIs on new connects. Neither
do they support any config settings to help us. You could rebuild your own chrome, but thats
not an option. So, you have to find a browser that is using system installed TLS libraries
to overlay the SSL functions we need.

`w3m` comes to the rescue!

```
~ $ LD_PRELOAD=/path/to/sniswitcher.so w3m https://www.bbc.co.uk
```

You will feel like 1990 (it could be worse!) but you have means to read news that your gov
think you shouldn't.

Now, part of the discussion can be the monopoly of browser vendors (or at least of browser engines)
and big-tech that always claims to help and free people but do very little on the technical side
(let alone the banning of domain-fronting) except making protocols overly complex to cast their
monoply in stone.

If this simple SNI switching trick does not help, we provide more tools to overcome blocking
such as [crash](https://github.com/stealth/crash) or [PSC](https://github.com/stealth/psc).

If you think that browsing in text mode anno 1990 is not sufficient to turn your outrage
into a revolution, you may use above mentioned `crash` in SNI mode to setup a full UDP/TCP
bridge for chrome. If outgoing connects are not possible anymore, client and server both
support active and passive connects on each side for a good reason.

