DROP TRIGGER notify_pricefeeds ON maker.price_feeds;

CREATE OR REPLACE FUNCTION notify_pricefeed() RETURNS trigger AS $$
BEGIN
  PERFORM pg_notify(
    CAST('postgraphile:price_feed' AS text),
    row_to_json(NEW)::text);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_pricefeeds
  AFTER INSERT ON maker.price_feeds
  FOR EACH ROW
  EXECUTE PROCEDURE notify_pricefeed();
