package booking

import "github.com/redis/go-redis/v9"

/*
 *  keys :
 *   1. booking:courtID:date
 *   2. bookings
 *  args:
 *   1. bookingId
 *   2. courtId
 *   3. startTime
 *   4. endTime
 *   5. guestName
 *
 */

var reserveScript = redis.NewScript(`
	local bookingsIndex = KEYS[1]
	local bookingsHashPrefix = KEYS[2]

	local bookingId = ARGV[1]
	local startTime = tonumber(ARGV[2])
	local endTime = tonumber(ARGV[3])
	local guestName = ARGV[4]

	local bookingHashKey = bookingsHashPrefix .. ":" .. bookingId

	-- Find existing bookings that start before endTime
	local candidates = redis.call('ZRANGEBYSCORE', bookingsIndex, '-inf', endTime)

	-- Check each for overlap
	for _, existingId in ipairs(candidates) do
		local existingKey = bookingsHashPrefix .. ":" .. existingId
		local existingEnd = redis.call('HGET', existingKey, 'end')

		if existingEnd and tonumber(existingEnd) > startTime then
			return 0
		end
	end

	-- No conflicts, create new booking
	redis.call('ZADD', bookingsIndex, startTime, bookingId)
	redis.call('HSET', bookingHashKey,
		'id', bookingId,
		'start', startTime,
		'end', endTime,
		'guestName', guestName
	)
	return 1
`)

var getBookingScript = redis.NewScript(`
	local bookingsHashPrefix = KEYS[1]
	local bookingId = ARGV[1]


	local bookingHashKey = bookingsHashPrefix .. ":" .. bookingId

	local booking = redis.call('HGETALL', bookingHashKey)
	return booking
	`)
