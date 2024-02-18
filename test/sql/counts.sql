SELECT 
    json_extract(data, "$.Host.Name") AS host,
    COUNT(*) AS count 
FROM jobRuns 
WHERE host LIKE 'srv-0%.acme.corp'
GROUP BY host ORDER BY count DESC LIMIT 100;