package db

const (
	SelectAllCountries = `select id from game_platform where release_year=1984;`

	AddPublisherQuery = `INSERT INTO publisher (publisher_name)
						 VALUES ((?));`

	DeletePublisherQuery = `DELETE FROM game_publisher WHERE publisher_id IN (SELECT id FROM publisher id WHERE publisher_name = (?));
							DELETE FROM publisher WHERE publisher_name = (?);`

	AddGamePublisherQuery = `INSERT INTO game_publisher (game_id, publisher_id)
							 VALUES ((SELECT id FROM game WHERE game_name = (?)), 
							 		 (SELECT id FROM publisher WHERE publisher_name = (?)));`

	DeleteGamePublisherQuery = `DELETE FROM region_sales WHERE game_platform_id IN (SELECT id from game_platform WHERE game_publisher_id = 
								(SELECT id FROM game_publisher WHERE (game_id = 
									(SELECT id from game WHERE game_name = (?))) AND publisher_id = 
									(SELECT id from publisher WHERE publisher_name = (?))));
	
								DELETE FROM game_platform WHERE game_publisher_id IN (SELECT id FROM game_publisher WHERE (game_id = 
									(SELECT id from game WHERE game_name = (?))) AND publisher_id = 
									(SELECT id from publisher WHERE publisher_name = (?)));
	
								DELETE FROM game_publisher WHERE id IN 
								(SELECT id FROM game_publisher WHERE (game_id = 
									(SELECT id from game WHERE game_name = (?))) AND publisher_id = 
									(SELECT id from publisher WHERE publisher_name = (?)));`

	DeleteGamePlatformByYearQuery = `DELETE FROM region_sales WHERE game_platform_id IN (SELECT id FROM game_platform WHERE release_year = (?));
									 DELETE FROM game_platform WHERE release_year = (?);`

	ChangePublisherQuery = `UPDATE publisher SET publisher_name = (?) WHERE publisher_name = (?);`
)
