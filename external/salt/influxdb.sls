influxdb:***REMOVED******REMOVED*** pillar.influxdb.version ***REMOVED******REMOVED***:
  dockerng.image_present

influxdb:
  dockerng.running:
    - image: influxdb:***REMOVED******REMOVED*** pillar.influxdb.version ***REMOVED******REMOVED***
    - network_mode: host
    - restart_policy: always
    - binds:
      - /var/lib/influxdb:/var/lib/influxdb
    - watch:
      - dockerng: influxdb:***REMOVED******REMOVED*** pillar.influxdb.version ***REMOVED******REMOVED***
