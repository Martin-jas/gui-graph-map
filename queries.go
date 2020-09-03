package main

const GET_COORDINATES = `
for c IN nodes
	FILTER c.Name == @name
	for coord in 1..1 inbound c graph 'worldGraphs'
    	let a = (
    	    FOR v, e, p IN 2..2 INBOUND coord GRAPH 'worldGraphs'
    	        RETURN v
    	    )
    	return {vertex: coord, owner:a}
`

const FIND_NODE_BY_NAME = `
FOR c IN nodes
	FILTER c.Name == @name
	RETURN c
`

const SOLVE_FIELDS = `
for c IN nodes
	FILTER c._id == @id
	FOR v, e IN 2..2 INBOUND c GRAPH 'worldGraphs'
		RETURN MERGE(v, {relation: PARSE_IDENTIFIER(e).collection})
`
