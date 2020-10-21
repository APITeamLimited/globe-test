/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a simple gRPC server that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It implements the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/reflection"
)

func init() ***REMOVED***
	rand.Seed(time.Now().UnixNano())
***REMOVED***

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type routeGuideServer struct ***REMOVED***
	UnimplementedRouteGuideServer
	savedFeatures []*Feature // read-only after initialized

	mu         sync.Mutex // protects routeNotes
	routeNotes map[string][]*RouteNote
***REMOVED***

// GetFeature returns the feature at the given point.
func (s *routeGuideServer) GetFeature(ctx context.Context, point *Point) (*Feature, error) ***REMOVED***

	n := rand.Intn(1000)
	time.Sleep(time.Duration(n) * time.Millisecond)

	for _, feature := range s.savedFeatures ***REMOVED***
		if proto.Equal(feature.Location, point) ***REMOVED***
			return feature, nil
		***REMOVED***
	***REMOVED***
	// No feature was found, return an unnamed feature
	return &Feature***REMOVED***Location: point***REMOVED***, nil
***REMOVED***

// ListFeatures lists all features contained within the given bounding Rectangle.
func (s *routeGuideServer) ListFeatures(rect *Rectangle, stream RouteGuide_ListFeaturesServer) error ***REMOVED***
	for _, feature := range s.savedFeatures ***REMOVED***
		if inRange(feature.Location, rect) ***REMOVED***
			time.Sleep(500 * time.Millisecond)
			if err := stream.Send(feature); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// RecordRoute records a route composited of a sequence of points.
//
// It gets a stream of points, and responds with statistics about the "trip":
// number of points,  number of known features visited, total distance traveled, and
// total time spent.
func (s *routeGuideServer) RecordRoute(stream RouteGuide_RecordRouteServer) error ***REMOVED***
	var pointCount, featureCount, distance int32
	var lastPoint *Point
	startTime := time.Now()
	for ***REMOVED***
		point, err := stream.Recv()
		if err == io.EOF ***REMOVED***
			endTime := time.Now()
			return stream.SendAndClose(&RouteSummary***REMOVED***
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			***REMOVED***)
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		pointCount++
		for _, feature := range s.savedFeatures ***REMOVED***
			if proto.Equal(feature.Location, point) ***REMOVED***
				featureCount++
			***REMOVED***
		***REMOVED***
		if lastPoint != nil ***REMOVED***
			distance += calcDistance(lastPoint, point)
		***REMOVED***
		lastPoint = point
	***REMOVED***
***REMOVED***

// RouteChat receives a stream of message/location pairs, and responds with a stream of all
// previous messages at each of those locations.
func (s *routeGuideServer) RouteChat(stream RouteGuide_RouteChatServer) error ***REMOVED***
	for ***REMOVED***
		in, err := stream.Recv()
		if err == io.EOF ***REMOVED***
			return nil
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		key := serialize(in.Location)

		s.mu.Lock()
		s.routeNotes[key] = append(s.routeNotes[key], in)
		// Note: this copy prevents blocking other clients while serving this one.
		// We don't need to do a deep copy, because elements in the slice are
		// insert-only and never modified.
		rn := make([]*RouteNote, len(s.routeNotes[key]))
		copy(rn, s.routeNotes[key])
		s.mu.Unlock()

		for _, note := range rn ***REMOVED***
			if err := stream.Send(note); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// loadFeatures loads features from a JSON file.
func (s *routeGuideServer) loadFeatures(filePath string) ***REMOVED***
	var data []byte
	if filePath != "" ***REMOVED***
		var err error
		data, err = ioutil.ReadFile(filePath)
		if err != nil ***REMOVED***
			log.Fatalf("Failed to load default features: %v", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		data = exampleData
	***REMOVED***
	if err := json.Unmarshal(data, &s.savedFeatures); err != nil ***REMOVED***
		log.Fatalf("Failed to load default features: %v", err)
	***REMOVED***
***REMOVED***

func toRadians(num float64) float64 ***REMOVED***
	return num * math.Pi / float64(180)
***REMOVED***

// calcDistance calculates the distance between two points using the "haversine" formula.
// The formula is based on http://mathforum.org/library/drmath/view/51879.html.
func calcDistance(p1 *Point, p2 *Point) int32 ***REMOVED***
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
***REMOVED***

func inRange(point *Point, rect *Rectangle) bool ***REMOVED***
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func serialize(point *Point) string ***REMOVED***
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
***REMOVED***

func newServer() *routeGuideServer ***REMOVED***
	s := &routeGuideServer***REMOVED***routeNotes: make(map[string][]*RouteNote)***REMOVED***
	s.loadFeatures(*jsonDBFile)
	return s
***REMOVED***

func main() ***REMOVED***
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil ***REMOVED***
		log.Fatalf("failed to listen: %v", err)
	***REMOVED***
	var opts []grpc.ServerOption
	if *tls ***REMOVED***
		if *certFile == "" ***REMOVED***
			*certFile = testdata.Path("server1.pem")
		***REMOVED***
		if *keyFile == "" ***REMOVED***
			*keyFile = testdata.Path("server1.key")
		***REMOVED***
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil ***REMOVED***
			log.Fatalf("Failed to generate credentials %v", err)
		***REMOVED***
		opts = []grpc.ServerOption***REMOVED***grpc.Creds(creds)***REMOVED***
	***REMOVED***
	grpcServer := grpc.NewServer(opts...)
	RegisterRouteGuideServer(grpcServer, newServer())
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)
***REMOVED***

// exampleData is a copy of testdata/route_guide_db.json. It's to avoid
// specifying file path with `go run`.
var exampleData = []byte(`[***REMOVED***
    "location": ***REMOVED***
        "latitude": 407838351,
        "longitude": -746143763
    ***REMOVED***,
    "name": "Patriots Path, Mendham, NJ 07945, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 408122808,
        "longitude": -743999179
    ***REMOVED***,
    "name": "101 New Jersey 10, Whippany, NJ 07981, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 413628156,
        "longitude": -749015468
    ***REMOVED***,
    "name": "U.S. 6, Shohola, PA 18458, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 419999544,
        "longitude": -740371136
    ***REMOVED***,
    "name": "5 Conners Road, Kingston, NY 12401, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 414008389,
        "longitude": -743951297
    ***REMOVED***,
    "name": "Mid Hudson Psychiatric Center, New Hampton, NY 10958, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 419611318,
        "longitude": -746524769
    ***REMOVED***,
    "name": "287 Flugertown Road, Livingston Manor, NY 12758, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 406109563,
        "longitude": -742186778
    ***REMOVED***,
    "name": "4001 Tremley Point Road, Linden, NJ 07036, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 416802456,
        "longitude": -742370183
    ***REMOVED***,
    "name": "352 South Mountain Road, Wallkill, NY 12589, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 412950425,
        "longitude": -741077389
    ***REMOVED***,
    "name": "Bailey Turn Road, Harriman, NY 10926, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 412144655,
        "longitude": -743949739
    ***REMOVED***,
    "name": "193-199 Wawayanda Road, Hewitt, NJ 07421, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 415736605,
        "longitude": -742847522
    ***REMOVED***,
    "name": "406-496 Ward Avenue, Pine Bush, NY 12566, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 413843930,
        "longitude": -740501726
    ***REMOVED***,
    "name": "162 Merrill Road, Highland Mills, NY 10930, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 410873075,
        "longitude": -744459023
    ***REMOVED***,
    "name": "Clinton Road, West Milford, NJ 07480, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 412346009,
        "longitude": -744026814
    ***REMOVED***,
    "name": "16 Old Brook Lane, Warwick, NY 10990, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 402948455,
        "longitude": -747903913
    ***REMOVED***,
    "name": "3 Drake Lane, Pennington, NJ 08534, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 406337092,
        "longitude": -740122226
    ***REMOVED***,
    "name": "6324 8th Avenue, Brooklyn, NY 11220, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 406421967,
        "longitude": -747727624
    ***REMOVED***,
    "name": "1 Merck Access Road, Whitehouse Station, NJ 08889, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 416318082,
        "longitude": -749677716
    ***REMOVED***,
    "name": "78-98 Schalck Road, Narrowsburg, NY 12764, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 415301720,
        "longitude": -748416257
    ***REMOVED***,
    "name": "282 Lakeview Drive Road, Highland Lake, NY 12743, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 402647019,
        "longitude": -747071791
    ***REMOVED***,
    "name": "330 Evelyn Avenue, Hamilton Township, NJ 08619, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 412567807,
        "longitude": -741058078
    ***REMOVED***,
    "name": "New York State Reference Route 987E, Southfields, NY 10975, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 416855156,
        "longitude": -744420597
    ***REMOVED***,
    "name": "103-271 Tempaloni Road, Ellenville, NY 12428, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404663628,
        "longitude": -744820157
    ***REMOVED***,
    "name": "1300 Airport Road, North Brunswick Township, NJ 08902, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 407113723,
        "longitude": -749746483
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 402133926,
        "longitude": -743613249
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 400273442,
        "longitude": -741220915
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 411236786,
        "longitude": -744070769
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 411633782,
        "longitude": -746784970
    ***REMOVED***,
    "name": "211-225 Plains Road, Augusta, NJ 07822, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 415830701,
        "longitude": -742952812
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 413447164,
        "longitude": -748712898
    ***REMOVED***,
    "name": "165 Pedersen Ridge Road, Milford, PA 18337, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 405047245,
        "longitude": -749800722
    ***REMOVED***,
    "name": "100-122 Locktown Road, Frenchtown, NJ 08825, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 418858923,
        "longitude": -746156790
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 417951888,
        "longitude": -748484944
    ***REMOVED***,
    "name": "650-652 Willi Hill Road, Swan Lake, NY 12783, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 407033786,
        "longitude": -743977337
    ***REMOVED***,
    "name": "26 East 3rd Street, New Providence, NJ 07974, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 417548014,
        "longitude": -740075041
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 410395868,
        "longitude": -744972325
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404615353,
        "longitude": -745129803
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 406589790,
        "longitude": -743560121
    ***REMOVED***,
    "name": "611 Lawrence Avenue, Westfield, NJ 07090, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 414653148,
        "longitude": -740477477
    ***REMOVED***,
    "name": "18 Lannis Avenue, New Windsor, NY 12553, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 405957808,
        "longitude": -743255336
    ***REMOVED***,
    "name": "82-104 Amherst Avenue, Colonia, NJ 07067, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 411733589,
        "longitude": -741648093
    ***REMOVED***,
    "name": "170 Seven Lakes Drive, Sloatsburg, NY 10974, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 412676291,
        "longitude": -742606606
    ***REMOVED***,
    "name": "1270 Lakes Road, Monroe, NY 10950, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 409224445,
        "longitude": -748286738
    ***REMOVED***,
    "name": "509-535 Alphano Road, Great Meadows, NJ 07838, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 406523420,
        "longitude": -742135517
    ***REMOVED***,
    "name": "652 Garden Street, Elizabeth, NJ 07202, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 401827388,
        "longitude": -740294537
    ***REMOVED***,
    "name": "349 Sea Spray Court, Neptune City, NJ 07753, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 410564152,
        "longitude": -743685054
    ***REMOVED***,
    "name": "13-17 Stanley Street, West Milford, NJ 07480, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 408472324,
        "longitude": -740726046
    ***REMOVED***,
    "name": "47 Industrial Avenue, Teterboro, NJ 07608, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 412452168,
        "longitude": -740214052
    ***REMOVED***,
    "name": "5 White Oak Lane, Stony Point, NY 10980, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 409146138,
        "longitude": -746188906
    ***REMOVED***,
    "name": "Berkshire Valley Management Area Trail, Jefferson, NJ, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404701380,
        "longitude": -744781745
    ***REMOVED***,
    "name": "1007 Jersey Avenue, New Brunswick, NJ 08901, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 409642566,
        "longitude": -746017679
    ***REMOVED***,
    "name": "6 East Emerald Isle Drive, Lake Hopatcong, NJ 07849, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 408031728,
        "longitude": -748645385
    ***REMOVED***,
    "name": "1358-1474 New Jersey 57, Port Murray, NJ 07865, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 413700272,
        "longitude": -742135189
    ***REMOVED***,
    "name": "367 Prospect Road, Chester, NY 10918, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404310607,
        "longitude": -740282632
    ***REMOVED***,
    "name": "10 Simon Lake Drive, Atlantic Highlands, NJ 07716, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 409319800,
        "longitude": -746201391
    ***REMOVED***,
    "name": "11 Ward Street, Mount Arlington, NJ 07856, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 406685311,
        "longitude": -742108603
    ***REMOVED***,
    "name": "300-398 Jefferson Avenue, Elizabeth, NJ 07201, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 419018117,
        "longitude": -749142781
    ***REMOVED***,
    "name": "43 Dreher Road, Roscoe, NY 12776, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 412856162,
        "longitude": -745148837
    ***REMOVED***,
    "name": "Swan Street, Pine Island, NY 10969, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 416560744,
        "longitude": -746721964
    ***REMOVED***,
    "name": "66 Pleasantview Avenue, Monticello, NY 12701, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 405314270,
        "longitude": -749836354
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 414219548,
        "longitude": -743327440
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 415534177,
        "longitude": -742900616
    ***REMOVED***,
    "name": "565 Winding Hills Road, Montgomery, NY 12549, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 406898530,
        "longitude": -749127080
    ***REMOVED***,
    "name": "231 Rocky Run Road, Glen Gardner, NJ 08826, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 407586880,
        "longitude": -741670168
    ***REMOVED***,
    "name": "100 Mount Pleasant Avenue, Newark, NJ 07104, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 400106455,
        "longitude": -742870190
    ***REMOVED***,
    "name": "517-521 Huntington Drive, Manchester Township, NJ 08759, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 400066188,
        "longitude": -746793294
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 418803880,
        "longitude": -744102673
    ***REMOVED***,
    "name": "40 Mountain Road, Napanoch, NY 12458, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 414204288,
        "longitude": -747895140
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 414777405,
        "longitude": -740615601
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 415464475,
        "longitude": -747175374
    ***REMOVED***,
    "name": "48 North Road, Forestburgh, NY 12777, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404062378,
        "longitude": -746376177
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 405688272,
        "longitude": -749285130
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 400342070,
        "longitude": -748788996
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 401809022,
        "longitude": -744157964
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404226644,
        "longitude": -740517141
    ***REMOVED***,
    "name": "9 Thompson Avenue, Leonardo, NJ 07737, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 410322033,
        "longitude": -747871659
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 407100674,
        "longitude": -747742727
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 418811433,
        "longitude": -741718005
    ***REMOVED***,
    "name": "213 Bush Road, Stone Ridge, NY 12484, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 415034302,
        "longitude": -743850945
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 411349992,
        "longitude": -743694161
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404839914,
        "longitude": -744759616
    ***REMOVED***,
    "name": "1-17 Bergen Court, New Brunswick, NJ 08901, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 414638017,
        "longitude": -745957854
    ***REMOVED***,
    "name": "35 Oakland Valley Road, Cuddebackville, NY 12729, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 412127800,
        "longitude": -740173578
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 401263460,
        "longitude": -747964303
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 412843391,
        "longitude": -749086026
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 418512773,
        "longitude": -743067823
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404318328,
        "longitude": -740835638
    ***REMOVED***,
    "name": "42-102 Main Street, Belford, NJ 07718, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 419020746,
        "longitude": -741172328
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404080723,
        "longitude": -746119569
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 401012643,
        "longitude": -744035134
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 404306372,
        "longitude": -741079661
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 403966326,
        "longitude": -748519297
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 405002031,
        "longitude": -748407866
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 409532885,
        "longitude": -742200683
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 416851321,
        "longitude": -742674555
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 406411633,
        "longitude": -741722051
    ***REMOVED***,
    "name": "3387 Richmond Terrace, Staten Island, NY 10303, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 413069058,
        "longitude": -744597778
    ***REMOVED***,
    "name": "261 Van Sickle Road, Goshen, NY 10924, USA"
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 418465462,
        "longitude": -746859398
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 411733222,
        "longitude": -744228360
    ***REMOVED***,
    "name": ""
***REMOVED***, ***REMOVED***
    "location": ***REMOVED***
        "latitude": 410248224,
        "longitude": -747127767
    ***REMOVED***,
    "name": "3 Hasta Way, Newton, NJ 07860, USA"
***REMOVED***]`)
