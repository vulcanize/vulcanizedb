ALTER TABLE maker.tend
 add constraint tend_bid_id_key unique (bid_id);

ALTER TABLE maker.dent
 add constraint dent_bid_id_key unique (bid_id);
