@heartbeat
Feature: heartbeat

  Scenario: Ok heartbeat
    When I GET http://localhost:8080/ping
    Then response status code should be {{status.Ok}}
