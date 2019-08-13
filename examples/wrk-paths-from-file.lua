

local _generated_requests = {}


function init (_arguments)
	for _index, _argument in ipairs (_arguments) do
		_generate_requests (_generated_requests, _argument)
	end
end


function request ()
	local _index = math.random (#_generated_requests)
	return _generated_requests[_index]
end


function _generate_requests (_requests, _path)
	local _tid = wrk.thread.tindex .. "/" .. wrk.thread.tcount
	print ("[ii] [" .. _tid .. "]  loading paths from `" .. _path .. "`...")
	local _index = 0
	local _wrk_path_prefix = wrk.path:gsub ("/$", "")
	for _wrk_path_suffix in io.lines (_path) do
		if math.fmod (_index, wrk.thread.tcount) == wrk.thread.tindex then
			if not _wrk_path_suffix:match ("/$") then
				_wrk_path_suffix = _wrk_path_suffix:gsub ("^/", "")
				_wrk_path = _wrk_path_prefix .. "/" .. _wrk_path_suffix
				_request = wrk.format (wrk.method, _wrk_path, wrk.headers, wrk.body)
				table.insert (_generated_requests, _request)
			end
		end
		_index = _index + 1
	end
	print ("[ii] [" .. _tid .. "]  loaded `" .. #_generated_requests .. '` of `' .. _index .. "` paths.")
end

