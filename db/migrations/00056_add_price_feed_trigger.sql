-- +goose Up
CREATE OR REPLACE FUNCTION notify_pricefeed() RETURNS trigger AS $$
BEGIN
  PERFORM pg_notify(
    CAST('postgraphile:price_feed' AS text),
    json_build_object('__node__', json_build_array('price_feeds', NEW.id))::text
  );
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_pricefeeds
  AFTER INSERT ON maker.price_feeds
  FOR EACH ROW
  EXECUTE PROCEDURE notify_pricefeed();


-- +goose Down
DROP TRIGGER notify_pricefeeds ON maker.price_feeds;
