/etc/telegraf.conf:
  file.managed:
    - source: salt://telegraf/telegraf.conf
    - template: jinja

/etc/telegraf.d:
  file.directory

telegraf:***REMOVED******REMOVED*** pillar.telegraf.version ***REMOVED******REMOVED***:
  dockerng.image_present

telegraf:
  dockerng.running:
    - image: telegraf:***REMOVED******REMOVED*** pillar.telegraf.version ***REMOVED******REMOVED***
    - cmd: -config /etc/telegraf.conf -config-directory /etc/telegraf.d
    - network_mode: host
    - restart_policy: always
    - binds:
      - /etc/telegraf.conf:/etc/telegraf.conf:ro
      - /etc/telegraf.d:/etc/telegraf.d:ro
    - watch:
      - file: /etc/telegraf.conf
      - dockerng: telegraf:***REMOVED******REMOVED*** pillar.telegraf.version ***REMOVED******REMOVED***
