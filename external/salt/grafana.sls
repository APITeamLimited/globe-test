grafana/grafana:***REMOVED******REMOVED*** pillar.grafana.version ***REMOVED******REMOVED***:
  dockerng.image_present

grafana:
  dockerng.running:
    - image: grafana/grafana:***REMOVED******REMOVED*** pillar.grafana.version ***REMOVED******REMOVED***
    - network_mode: host
    - restart_policy: always
    - environment:
      - GF_AUTH_ANONYMOUS_ENABLED: "True"
      - GF_AUTH_ANONYMOUS_ORG_ROLE: Admin
    - binds:
      - /var/lib/grafana:/var/lib/grafana
    - watch:
      - dockerng: grafana/grafana:***REMOVED******REMOVED*** pillar.grafana.version ***REMOVED******REMOVED***
