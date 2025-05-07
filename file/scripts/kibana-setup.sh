#!/bin/bash

# Wait for Elasticsearch to be available
echo "Waiting for Elasticsearch to be available..."
until curl -s http://host.docker.internal:9200/_cluster/health | grep -E '"status":"(green|yellow)"' > /dev/null; do
  sleep 2
done
echo "Elasticsearch is up and running with green or yellow status!"

# Create index pattern
echo "Creating index pattern for file-service logs..."
curl -X POST "http://host.docker.internal:5601/api/saved_objects/index-pattern/file-service-logs-*" \
-H 'kbn-xsrf: true' \
-H 'Content-Type: application/json' \
-d '
{
  "attributes": {
    "title": "file-service-logs-*",
    "timeFieldName": "@timestamp"
  }
}'

echo "Creating visualization for log levels..."
curl -X POST "http://host.docker.internal:5601/api/saved_objects/visualization/file-service-log-levels" \
-H 'kbn-xsrf: true' \
-H 'Content-Type: application/json' \
-d '
{
  "attributes": {
    "title": "File Service Log Levels",
    "visState": "{\"title\":\"File Service Log Levels\",\"type\":\"pie\",\"params\":{\"type\":\"pie\",\"addTooltip\":true,\"addLegend\":true,\"legendPosition\":\"right\",\"isDonut\":false},\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"count\",\"schema\":\"metric\",\"params\":{}},{\"id\":\"2\",\"enabled\":true,\"type\":\"terms\",\"schema\":\"segment\",\"params\":{\"field\":\"level.keyword\",\"size\":5,\"order\":\"desc\",\"orderBy\":\"1\",\"otherBucket\":false,\"otherBucketLabel\":\"Other\",\"missingBucket\":false,\"missingBucketLabel\":\"Missing\"}}]}",
    "uiStateJSON": "{}",
    "description": "Pie chart showing distribution of log levels",
    "version": 1,
    "kibanaSavedObjectMeta": {
      "searchSourceJSON": "{\"index\":\"file-service-logs-*\",\"query\":{\"query\":\"\",\"language\":\"kuery\"},\"filter\":[]}"
    }
  }
}'

echo "Creating visualization for file operations..."
curl -X POST "http://host.docker.internal:5601/api/saved_objects/visualization/file-service-operations" \
-H 'kbn-xsrf: true' \
-H 'Content-Type: application/json' \
-d '
{
  "attributes": {
    "title": "File Service Operations",
    "visState": "{\"title\":\"File Service Operations\",\"type\":\"histogram\",\"params\":{\"type\":\"histogram\",\"grid\":{\"categoryLines\":false},\"categoryAxes\":[{\"id\":\"CategoryAxis-1\",\"type\":\"category\",\"position\":\"bottom\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\"},\"labels\":{\"show\":true,\"truncate\":100},\"title\":{}}],\"valueAxes\":[{\"id\":\"ValueAxis-1\",\"name\":\"LeftAxis-1\",\"type\":\"value\",\"position\":\"left\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\",\"mode\":\"normal\"},\"labels\":{\"show\":true,\"rotate\":0,\"filter\":false,\"truncate\":100},\"title\":{\"text\":\"Count\"}}],\"seriesParams\":[{\"show\":true,\"type\":\"histogram\",\"mode\":\"stacked\",\"data\":{\"label\":\"Count\",\"id\":\"1\"},\"valueAxis\":\"ValueAxis-1\",\"drawLinesBetweenPoints\":true,\"lineWidth\":2,\"showCircles\":true}],\"addTooltip\":true,\"addLegend\":true,\"legendPosition\":\"right\",\"times\":[],\"addTimeMarker\":false,\"labels\":{\"show\":false},\"thresholdLine\":{\"show\":false,\"value\":10,\"width\":1,\"style\":\"full\",\"color\":\"#E7664C\"},\"dimensions\":{\"x\":{\"accessor\":0,\"format\":{\"id\":\"terms\",\"params\":{\"id\":\"string\",\"otherBucketLabel\":\"Other\",\"missingBucketLabel\":\"Missing\"}},\"params\":{},\"label\":\"message.keyword: Descending\",\"aggType\":\"terms\"},\"y\":[{\"accessor\":1,\"format\":{\"id\":\"number\"},\"params\":{},\"label\":\"Count\",\"aggType\":\"count\"}],\"series\":[{\"accessor\":2,\"format\":{\"id\":\"terms\",\"params\":{\"id\":\"string\",\"otherBucketLabel\":\"Other\",\"missingBucketLabel\":\"Missing\"}},\"params\":{},\"label\":\"level.keyword: Descending\",\"aggType\":\"terms\"}]}},\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"count\",\"schema\":\"metric\",\"params\":{}},{\"id\":\"2\",\"enabled\":true,\"type\":\"terms\",\"schema\":\"segment\",\"params\":{\"field\":\"message.keyword\",\"orderBy\":\"1\",\"order\":\"desc\",\"size\":10,\"otherBucket\":false,\"otherBucketLabel\":\"Other\",\"missingBucket\":false,\"missingBucketLabel\":\"Missing\",\"include\":\"(Upload|Retrieving|Deleting|Updating)\"}},{\"id\":\"3\",\"enabled\":true,\"type\":\"terms\",\"schema\":\"group\",\"params\":{\"field\":\"level.keyword\",\"orderBy\":\"1\",\"order\":\"desc\",\"size\":5,\"otherBucket\":false,\"otherBucketLabel\":\"Other\",\"missingBucket\":false,\"missingBucketLabel\":\"Missing\"}}]}",
    "uiStateJSON": "{}",
    "description": "Histogram of file service operations",
    "version": 1,
    "kibanaSavedObjectMeta": {
      "searchSourceJSON": "{\"index\":\"file-service-logs-*\",\"query\":{\"query\":\"\",\"language\":\"kuery\"},\"filter\":[]}"
    }
  }
}'

echo "Creating dashboard..."
curl -X POST "http://host.docker.internal:5601/api/saved_objects/dashboard/file-service-monitoring" \
-H 'kbn-xsrf: true' \
-H 'Content-Type: application/json' \
-d '
{
  "attributes": {
    "title": "File Service Monitoring",
    "hits": 0,
    "description": "Dashboard for monitoring file service operations",
    "panelsJSON": "[{\"gridData\":{\"x\":0,\"y\":0,\"w\":24,\"h\":15,\"i\":\"1\"},\"version\":\"7.10.2\",\"panelIndex\":\"1\",\"embeddableConfig\":{\"title\":\"File Service Log Levels\"},\"panelRefName\":\"panel_0\"},{\"gridData\":{\"x\":24,\"y\":0,\"w\":24,\"h\":15,\"i\":\"2\"},\"version\":\"7.10.2\",\"panelIndex\":\"2\",\"embeddableConfig\":{\"title\":\"File Service Operations\"},\"panelRefName\":\"panel_1\"}]",
    "optionsJSON": "{\"hidePanelTitles\":false,\"useMargins\":true}",
    "version": 1,
    "timeRestore": false,
    "kibanaSavedObjectMeta": {
      "searchSourceJSON": "{\"query\":{\"language\":\"kuery\",\"query\":\"\"},\"filter\":[]}"
    }
  },
  "references": [
    {
      "name": "panel_0",
      "type": "visualization",
      "id": "file-service-log-levels"
    },
    {
      "name": "panel_1",
      "type": "visualization",
      "id": "file-service-operations"
    }
  ]
}'

echo "Kibana setup completed! Dashboard and visualizations have been created." 