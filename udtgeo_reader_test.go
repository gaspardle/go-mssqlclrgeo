package mssqlclrgeo

import (
	"encoding/hex"
	"testing"
)

func TestGeographt(t *testing.T) {
	values := []struct {
		hex       string
		nbPoints  int
		nbShapes  int
		typeShape SHAPE
		srid      int32
	}{
		//POINT (20 40)', 4326
		{"E6100000010C00000000000044400000000000003440",
			1, 1, SHAPE_POINT, 4326},
		//'LINESTRING (21 45, -45 20), 4326
		{"E6100000011400000000008046400000000000003540000000000000344000000000008046C0",
			2, 1, SHAPE_LINESTRING, 4326},
		//geom LINESTRING (100 100, 20 180, 180 180)
		{"E610000001040300000000000000000059400000000000005940000000000000344000000000008066400000000000806640000000000080664001000000010000000001000000FFFFFFFF0000000002",
			3, 1, SHAPE_LINESTRING, 4326},
		//POLYGON ((74.1160 -10.1953,73.8248 -44.1210,82.1903 -44.1210,82.2616 -10.0195,74.1160 -10.1953)) 4326
		{"E61000000104050000004ED1915CFE6324C08195438B6C875240736891ED7C0F46C0696FF085C9745240736891ED7C0F46C0711B0DE02D8C5440448B6CE7FB0924C08D28ED0DBE9054404ED1915CFE6324C08195438B6C87524001000000020000000001000000FFFFFFFF0000000003",
			5, 1, SHAPE_POLYGON, 4326},
		//MULTIPOLYGON (((40 40, 20 45, 45 30, 40 40)), ((20 35, 45 20, 30 5, 10 10, 10 30, 20 35), (30 20, 20 25, 20 15, 30 20))) 4326
		{"E610000001040E000000000000000000444000000000000044400000000000003440000000000080464000000000008046400000000000003E400000000000004440000000000000444000000000000034400000000000804140000000000080464000000000000034400000000000003E4000000000000014400000000000002440000000000000244000000000000024400000000000003E40000000000000344000000000008041400000000000003E4000000000000034400000000000003440000000000000394000000000000034400000000000002E400000000000003E4000000000000034400300000002000000000204000000000A00000003000000FFFFFFFF0000000006000000000000000003000000000100000003",
			14, 3, SHAPE_MULTIPOLYGON, 4326},
		//GEOMETRYCOLLECTION(POINT(4 6),LINESTRING(4 6,7 10)) 4326
		{"E6100000010403000000000000000000184000000000000010400000000000001840000000000000104000000000000024400000000000001C40020000000100000000010100000003000000FFFFFFFF0000000007000000000000000001000000000100000002",
			3, 3, SHAPE_GEOMETRY_COLLECTION, 4326},
		//CURVEPOLYGON(COMPOUNDCURVE((0 0, 0 2, 2 2), CIRCULARSTRING (2 2, 1 0, 0 0))) 4326
		{"E61000000224050000000000000000000000000000000000000000000000000000400000000000000000000000000000004000000000000000400000000000000000000000000000F03F0000000000000000000000000000000001000000030000000001000000FFFFFFFF000000000A03000000020003",
			5, 1, SHAPE_CURVE_POLYGON, 4326},
	}

	for _, v := range values {
		udtbin, _ := hex.DecodeString(v.hex)
		g, err := ReadGeography(udtbin)
		if err != nil {
			t.Error(err)
		}

		if len(g.Points) != v.nbPoints {
			t.Errorf("(type#%d) points count doesn't match Value: %d Expected: %d", v.typeShape, len(g.Points), v.nbPoints)
		}
		if len(g.Shapes) != v.nbShapes {
			t.Errorf("(type#%d) shape count doesn't match, Value: %d Expected: %d", v.typeShape, len(g.Shapes), v.nbShapes)
		}
		if g.SRID != v.srid {
			t.Errorf("(type#%d) Srid doesn't match, Value: %d Expected: %d", v.typeShape, g.SRID, v.srid)
		}
	}
}
func TestInvalidData(t *testing.T) {
	// invalid version 153
	udtbin, _ := hex.DecodeString("0F270000990C00000000000000000000000000000000")
	_, err := ReadGeometry(udtbin)
	if err == nil {
		t.Error(" should have returned an error")
	}

	//  P and L are set
	udtbin, _ = hex.DecodeString("0F27000001FF00000000000000000000000000000000")
	_, err = ReadGeometry(udtbin)
	if err == nil {
		t.Error(" should have returned an error")
	}
}

func TestInvalidGeography(t *testing.T) {
	// srid = 0
	udtbin, _ := hex.DecodeString("00000000010C00000000000034400000000000004440")
	_, err := ReadGeography(udtbin)
	if err == nil {
		t.Error("geography with srid 0 should have returned an error")
	}
	// invalid srid
	udtbin, _ = hex.DecodeString("D2040000010C569FABADD8074E40C5FEB27BF2B02840")
	g, err := ReadGeography(udtbin)
	if err == nil {
		t.Errorf("geography with srid %d should have returned an error", g.SRID)
	}

	//'POINT (32000 -45) 4326
	udtbin, _ = hex.DecodeString("E6100000010C000000000040DF4000000000008046C0")
	g, err = ReadGeography(udtbin)
	if err == nil {
		t.Error("longitude not between -15069 and 15069 should fail")
	}

	//'POINT (12.3 -91) 4326
	udtbin, _ = hex.DecodeString("E6100000010C9A999999999928400000000000C056C0")
	g, err = ReadGeography(udtbin)
	if err == nil {
		t.Error("latitude not between -90 and 90 should fail")
	}
}

func TestGeometry(t *testing.T) {
	//POINT (20 40)', 0
	//{"00000000010C00000000000034400000000000004440"
	//LINESTRING (100 100, 20 180, 180 180) 4326
	//{"E610000001040300000000000000000059400000000000005940000000000000344000000000008066400000000000806640000000000080664001000000010000000001000000FFFFFFFF0000000002"
	//CURVEPOLYGON(COMPOUNDCURVE((0 0, 0 2, 2 2), CIRCULARSTRING (2 2, 1 0, 0 0)))
	//"00000000020405000000000000000000000000000000000000000000000000000000000000000000004000000000000000400000000000000040000000000000F03F00000000000000000000000000000000000000000000000001000000030000000001000000FFFFFFFF000000000A03000000020003"

}
