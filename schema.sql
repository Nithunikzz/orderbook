-- Currency pairs table
CREATE TABLE IF NOT EXISTS currency_pairs (
    id SERIAL PRIMARY KEY,
    base TEXT NOT NULL,
    quote TEXT NOT NULL
);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    pair TEXT NOT NULL,
    side TEXT CHECK (side IN ('buy', 'sell')) NOT NULL,
    price DECIMAL NOT NULL,
    quantity DECIMAL NOT NULL,
    filled_quantity DECIMAL DEFAULT 0,
    status TEXT CHECK (status IN ('open', 'filled', 'partial', 'cancelled')) NOT NULL,
    fee DECIMAL DEFAULT 0,
    created_at TIMESTAMP DEFAULT now()
    
);
