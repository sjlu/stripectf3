package com.stripe.ctf.instantcodesearch

import java.io._
import java.nio.file._

import com.twitter.concurrent.Broker
import scala.collection.immutable.HashMap

abstract class SearchResult
case class Match(path: String, line: Int) extends SearchResult
case class Done() extends SearchResult

class Searcher(indexPath : String)  {

  def search(needle : String, b : Broker[SearchResult], indexer : Indexer) = {
    for (path <- indexer.indexMap.keys) {
      for (m <- tryPath(path, needle, indexer)) {
        b !! m
      }
    }

    b !! new Done()
  }

  def tryPath(path: String, needle: String, indexer : Indexer) : Iterable[SearchResult] = {
    try {
      val text : String = indexer.indexMap(path)
      if (text.contains(needle)) {
        var line = 0
        return text.split("\n").zipWithIndex.
          filter { case (l,n) => l.contains(needle) }.
          map { case (l,n) => new Match(path, n+1) }
      }
    } catch {
      case e: IOException => {
        return Nil
      }
    }

    return Nil
  }

  def readIndex(path: String) : Index = {
    new ObjectInputStream(new FileInputStream(new File(path))).readObject.asInstanceOf[Index]
  }
}
