package booking

import "github.com/redis/go-redis/v9"

/*
 *  keys :
 *   1. booking:courtID:date
 *   2. bookings
 *  args:
 *   1. bookingId
 *   2. startTime
 *   3. endTime
 *   4. guestName
 */

var reserveScript = redis.NewScript(`
	local bookingsIndex = KEYS[1]
	local bookingsHashPrefix = KEYS[2]

	local _, _, courtID = string.find(bookingsIndex, "bookings:{(.+)}")

	local bookingId = ARGV[1]
	local startTime = tonumber(ARGV[2])
	local endTime = tonumber(ARGV[3])
	local guestName = ARGV[4]

	local bookingHashKey = bookingsHashPrefix .. ":" .. bookingId

	-- Find existing bookings that start before endTime
	local candidates = redis.call('ZRANGEBYSCORE', bookingsIndex, '-inf', endTime)

	for _, existingId in ipairs(candidates) do
		local existingKey = bookingsHashPrefix .. ":" .. existingId
		local existingEnd = redis.call('HGET', existingKey, 'end')

		if existingEnd and tonumber(existingEnd) > startTime then
			return 0
		end
	end

	-- No conflict, create new booking
	redis.call('ZADD', bookingsIndex, startTime, bookingId)
	redis.call('HSET', bookingHashKey, 'id', bookingId, 'start', startTime, 'end', endTime, 'guestName', guestName, 'courtId', courtID)
	return 1
`)

/*
 *   1. booking:bookingId
 */
var getBookingScript = redis.NewScript(`
	local bookingHashKey = KEYS[1]

	local booking = redis.call('HGETALL', bookingHashKey)
	return booking
	`)

var releaseBookingScript = redis.NewScript(`
	local bookingsHashKey = KEYS[1]

	local bookingID = ARGV[1]

	local bookingCourtID = redis.call('HGET', bookingsHashKey, 'courtId')

	if not bookingCourtID then
		return 0
	end

	local bookingsIndex = "bookings:{" .. bookingCourtID .. "}"
	redis.call('ZREM', bookingsIndex, bookingID)
	redis.call('DEL', bookingsHashKey)

	return 1
	`)
