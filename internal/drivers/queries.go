package drivers

const (
	queryCreateQuote = `
	INSERT INTO quotes (id, author, text)
	VALUES ($1, $2, $3)
`
	queryDeleteQuote = `
	DELETE FROM quotes 
	WHERE id = $1
`
	queryGetAllQuotes = `
	SELECT id, author, text
	FROM quotes
`
	queryGetQuoteByAuthor = `
	SELECT id, text
	FROM quotes
	WHERE author = $1
`
	queryGetRandomQuote = `
	SELECT id, author, text
	FROM quotes 
	ORDER BY RANDOM()
	LIMIT 1
`
	queryGetQuoteById = `
	SELECT author, text
	FROM quotes
	WHERE id = $1
`
)
